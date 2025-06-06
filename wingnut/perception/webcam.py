import cv2
import asyncio
import logging
from config import CAMERA_INDEX, CASCADE_PATH

logger = logging.getLogger(__name__)

class WebcamFaceRecognition:
    def __init__(self, camera_index=CAMERA_INDEX, cascade_path=CASCADE_PATH):
        self.camera_index = camera_index
        self.face_cascade = cv2.CascadeClassifier(cascade_path)
        self._running = False

    async def start(self):
        """Continuously grab frames & detect faces in a background thread."""
        logger.info("Starting webcam face recognition")
        self._running = True
        loop = asyncio.get_running_loop()
        # run the blocking capture loop in a thread pool
        faces = await loop.run_in_executor(None, self._capture_loop)
        return faces

    def _capture_loop(self):
        cap = cv2.VideoCapture(self.camera_index)
        while self._running:
            ret, frame = cap.read()
            if not ret:
                logger.error("Failed to read frame from camera")
                break
            gray   = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
            faces  = self.face_cascade.detectMultiScale(gray, 1.1, 4)
            logger.debug(f"Detected {len(faces)} faces")
            return faces
        cap.release()

    async def stop(self):
        logger.info("Stopping webcam")
        self._running = False
