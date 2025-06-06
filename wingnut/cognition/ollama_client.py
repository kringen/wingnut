import aiohttp
import asyncio
import logging
import json
from config import OLLAMA_ENDPOINT, OLLAMA_MODEL, SPEAK_ENABLED

if SPEAK_ENABLED:
    try:
        import pyttsx3
    except ImportError:
        SPEAK_ENABLED = False
        logging.warning("pyttsx3 not installed, speech output will be disabled.")

logger = logging.getLogger(__name__)

class OllamaClient:
    def __init__(self, endpoint=OLLAMA_ENDPOINT):
        self.endpoint = endpoint
        self._session = None

    async def _get_session(self):
        if not self._session:
            self._session = aiohttp.ClientSession()
        return self._session

    async def ask(self, question: str):
        logger.info(f"Asking Ollama: {question!r}")
        sess = await self._get_session()
        payload = {"model": OLLAMA_MODEL, "prompt": question, "stream": True}
        logger.debug(f"Ollama request payload: {payload!r}")
        responses = []
        if SPEAK_ENABLED:
            engine = pyttsx3.init()
            logger.info("Speech output is enabled, initializing text-to-speech engine.")
        else:
            logger.warning("Speech output is disabled. Set SPEAK_ENABLED to True to enable.")
        
        async with sess.post(self.endpoint, json=payload) as resp:
            async for line in resp.content:
                if line:
                    line_json = json.loads(line)
                    if line_json.get("response"):
                        resp_chunk = line_json["response"]
                        responses.append(resp_chunk)
                        if SPEAK_ENABLED:
                            # Use text-to-speech to read the response chunk
                            engine.say(resp_chunk)
                            engine.runAndWait()
                        print(resp_chunk)
                    if line_json.get("done"):
                        if line_json["done"] == True:
                            print("Ollama stream ended")
                            await sess.close()
                            break
        
        logger.info("Ollama interaction completed")   
        return ''.join(responses)

    async def close(self):
        if self._session:
            await self._session.close()