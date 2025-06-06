import asyncio
import logging

from config import *
from perception.webcam import WebcamFaceRecognition
from actuation.drive    import DriveFunctions
from actuation.servo    import ServoFunctions
from cognition.ollama_client import OllamaClient
from cognition.keyword_listener import KeywordListener
from core.wingnut import WingnutRobot

logger = logging.getLogger(__name__)

async def main():
    # hardware placeholders — wire your real motor/servo objects here
    left_motor  = ...
    right_motor = ...
    head_servo  = ...

    webcam   = WebcamFaceRecognition()
    drive    = DriveFunctions(left_motor, right_motor)
    servo    = ServoFunctions(head_servo)
    ollama   = OllamaClient()
    listener = KeywordListener(callback_coro=ollama.ask)

    robot = WingnutRobot(webcam, drive, servo, ollama, listener)
    # try:
    #     response = await robot.ask_ollama("What is the meaning of life?")
    #     print("Ollama response:", response)
    # except Exception as e:
    #     logger.error(f"Error during Ollama interaction: {e}")
    try:
        faces = await robot.start()
        print("Detected faces:", faces)
    finally:
        await robot.stop()

if __name__ == '__main__':
    # catch top‐level errors
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Interrupted by user")
    except Exception:
        logger.exception("Unhandled exception")
