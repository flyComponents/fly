import requests
import time
import uuid

BASE_URL = "http://localhost:8080"

# можно задать вручную или генерить
AGENT_ID = str(uuid.uuid4())

HEARTBEAT_INTERVAL = 2  # секунды


def register_agent(agent_id):
    url = f"{BASE_URL}/agents/register"
    
    payload = {
        "id": agent_id
    }

    try:
        response = requests.post(url, json=payload)
        response.raise_for_status()
        print(f"[INFO] Registered agent: {agent_id}")
        return agent_id
    except requests.RequestException as e:
        print(f"[ERROR] Registration failed: {e}")
        return None


def send_heartbeat(agent_id):
    url = f"{BASE_URL}/agents/{agent_id}/heartbeat"
    
    try:
        response = requests.post(url)
        response.raise_for_status()
        print(f"[HEARTBEAT] sent for {agent_id}")
    except requests.RequestException as e:
        print(f"[ERROR] Heartbeat failed: {e}")


def main():
    agent_id = register_agent(AGENT_ID)
    
    if not agent_id:
        return

    while True:
        send_heartbeat(agent_id)
        time.sleep(HEARTBEAT_INTERVAL)


if __name__ == "__main__":
    main()