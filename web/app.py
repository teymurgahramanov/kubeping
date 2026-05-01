from flask import Flask, render_template, request, session, redirect, url_for
from kubernetes import client, config as k8s_config
from config import config
import concurrent.futures
import os
import requests

k8s_config.load_incluster_config()
with open("/var/run/secrets/kubernetes.io/serviceaccount/namespace") as f:
    current_namespace = f.read().strip()
v1 = client.CoreV1Api()

label_selector = "kubeping/component=exporter"
exporter_probe_path = '/probe'
app_version = config.APP_VERSION
public_ip_url = config.PUBLIC_IP_URL

app = Flask(__name__)
app.secret_key = os.environ.get(
    "SECRET_KEY",
    "insecure-development-key-set-secret-key-env-in-production",
)

def get_public_ip():
    try:
        response = requests.get(public_ip_url, headers={"Accept": "text/plain"}, timeout=5)
        response.raise_for_status()
        return {"ip": response.text.strip(), "source": public_ip_url}
    except Exception as e:
        return {"ip": "N/A"}

@app.route('/healthz')
def healthz():
    return 'ok', 200

@app.route('/')
def index():
    public_ip = get_public_ip()
    results = session.pop('results', None)
    return render_template('index.html', version=app_version, public_ip=public_ip, results=results)

@app.route('/ping', methods=['POST'])
def ping():
    data = {
        "module": "tcp",
        "address": request.form['address'],
        "timeout": int(request.form['timeout'])
    }
    exporters = {}
    session['results'] = []
    pods = v1.list_namespaced_pod(namespace=current_namespace, label_selector=label_selector)
    
    if not pods.items:
        session['results'].append({
            "host": 0,
            "result": f"Can't find kubeping-exporter pods with label selector {label_selector}"
        })
    else:
        for pod in pods.items:
            container_port = pod.spec.containers[0].ports[0].container_port
            exporters[pod.metadata.name] = {
                "host": pod.status.host_ip,
                "api_url": f"http://{pod.status.pod_ip}:{container_port}{exporter_probe_path}"
            }

        with concurrent.futures.ThreadPoolExecutor() as executor:
            futures = [executor.submit(probe, exporter, data) for exporter in exporters.values()]
            for future in concurrent.futures.as_completed(futures):
                result = future.result()
                if result:
                    session['results'].append(result)

    return redirect(url_for('index'))

def probe(exporter, data):
    try:
        response = requests.post(exporter["api_url"], json=data, timeout=data['timeout']*2)
        response_data = response.json()
        result = response_data.get("error", response_data.get("result"))
    except Exception as e:
        result = str(e)

    return {
        "host": exporter["host"],
        "result": result
    }

if __name__ == '__main__':
    app.run()