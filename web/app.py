from flask import Flask, render_template, request, session, redirect, url_for
from kubernetes import client, config as k8s_config
from config import config
import requests
import concurrent.futures

k8s_config.load_incluster_config()
with open("/var/run/secrets/kubernetes.io/serviceaccount/namespace") as f:
    current_namespace = f.read().strip()
v1 = client.CoreV1Api()

label_selector = "kubeping/component=exporter"
exporter_port = 8000
exporter_probe_path = '/probe'
app_version = config.APP_VERSION

app = Flask(__name__)
app.secret_key = 'secret'

@app.route('/')
def index():
    return render_template('index.html', version=app_version)

@app.route('/submit', methods=['POST'])
def submit():
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
            exporters[pod.metadata.name] = {
                "host": pod.status.host_ip,
                "api_url": f"http://{pod.status.pod_ip}:{exporter_port}{exporter_probe_path}"
            }

        with concurrent.futures.ThreadPoolExecutor() as executor:
            futures = [executor.submit(probe, exporter, data) for exporter in exporters.values()]
            for future in concurrent.futures.as_completed(futures):
                result = future.result()
                if result:
                    session['results'].append(result)

    return render_template('index.html', version=app_version, results=session['results'])

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