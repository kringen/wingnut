import logging

# hardware & service constants
CAMERA_INDEX   = 0
CASCADE_PATH   = 'haarcascades/haarcascade_frontalface_default.xml'
OLLAMA_ENDPOINT = 'http://192.168.1.66:11434/api/generate'
OLLAMA_MODEL    = 'llama3.2:latest'
SPEAK_ENABLED   = False
KEYWORD = 'hello'
# logging setup (you can tweak level & format here)
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s %(name)s %(levelname)s: %(message)s'
)
