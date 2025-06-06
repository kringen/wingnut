import asyncio
import logging

logger = logging.getLogger(__name__)

class ServoFunctions:
    def __init__(self, head_servo):
        self.head_servo = head_servo

    async def move_head_left(self):
        logger.info("Moving head left")
        await asyncio.to_thread(self.head_servo.set_angle, -45)

    async def move_head_right(self):
        logger.info("Moving head right")
        await asyncio.to_thread(self.head_servo.set_angle, 45)
