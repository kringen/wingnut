import asyncio
import logging
from perception.webcam import WebcamFaceRecognition
from actuation.drive    import DriveFunctions
from actuation.servo    import ServoFunctions
from cognition.ollama_client import OllamaClient
from cognition.keyword_listener import KeywordListener

logger = logging.getLogger(__name__)

class WingnutRobot:
    def __init__(self,
                 webcam: WebcamFaceRecognition,
                 drive: DriveFunctions,
                 servo: ServoFunctions,
                 ollama: OllamaClient,
                 listener: KeywordListener):
        self.webcam   = webcam
        self.drive    = drive
        self.servo    = servo
        self.ollama   = ollama
        self.listener = listener

    async def start(self):
        logger.info("Wingnut starting all coroutines")
        # start listener & webcam in parallel
        tasks = [
            asyncio.create_task(self.listener.run(), name="keyword-listener"),
            asyncio.create_task(self.webcam.start(), name="face-recognition"),
        ]
        # wait until one task returns (e.g. faces detected) or is cancelled
        done, pending = await asyncio.wait(
            tasks,
            return_when=asyncio.FIRST_COMPLETED
        )

        # grab result from face-recognition
        for t in done:
            if t.get_name() == "face-recognition":
                faces = t.result()
                logger.info(f"Faces: {faces}")
                # you could route to drive/servo here
        # cancel leftovers
        for t in pending:
            t.cancel()
        return faces

    async def stop(self):
        logger.info("Wingnut shutting down")
        await self.webcam.stop()
        await self.listener.stop()
        await self.ollama.close()
        
    async def ask_ollama(self, question: str):
        logger.info(f"Asking Ollama: {question!r}")
        response = await self.ollama.ask(question)
        logger.info(f"Ollama response: {response}")
        return response