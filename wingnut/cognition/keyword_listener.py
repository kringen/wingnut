import sounddevice as sd
import numpy as np
import asyncio
import logging
from config import KEYWORD

logger = logging.getLogger(__name__)

class KeywordListener:
    def __init__(self, callback_coro, samplerate=44100, blocksize=1024):
        """
        callback_coro: async function to call when keyword is detected
        samplerate: Audio sample rate
        blocksize: Number of frames per buffer
        """
        self.callback = callback_coro
        self.samplerate = samplerate
        self.blocksize = blocksize
        self._running = False

    async def run(self):
        logger.info("Keyword listener starting")
        self._running = True

        def callback(indata, frames, time, status):
            if status:
                logger.warning(f"Audio stream error: {status}")
            if b'{ KEYWORD }' in indata.tobytes():  # Simulated keyword detection
                logger.info("Keyword detected")
                asyncio.create_task(self.callback("What is your name?"))

        with sd.InputStream(samplerate=self.samplerate, blocksize=self.blocksize,
                            channels=1, dtype='int16', callback=callback):
            while self._running:
                await asyncio.sleep(0.1)  # Keep loop alive

    async def stop(self):
        logger.info("Stopping keyword listener")
        self._running = False
