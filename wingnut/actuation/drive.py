import asyncio
import logging

logger = logging.getLogger(__name__)

class DriveFunctions:
    def __init__(self, left_motor, right_motor):
        self.left_motor  = left_motor
        self.right_motor = right_motor

    async def move_left(self, speed, direction):
        logger.debug(f"move_left speed={speed} dir={direction}")
        if direction == 'forward':
            await asyncio.to_thread(self.left_motor.set_speed, speed)
            await asyncio.to_thread(self.right_motor.set_speed, -speed)
        else:  # backward
            await asyncio.to_thread(self.left_motor.set_speed, -speed)
            await asyncio.to_thread(self.right_motor.set_speed, speed)

    async def move_right(self, speed, direction):
        logger.debug(f"move_right speed={speed} dir={direction}")
        if direction == 'forward':
            await asyncio.to_thread(self.left_motor.set_speed, speed)
            await asyncio.to_thread(self.right_motor.set_speed, speed)
        else:
            await asyncio.to_thread(self.left_motor.set_speed, -speed)
            await asyncio.to_thread(self.right_motor.set_speed, -speed)

    async def go_forward(self, speed):
        logger.info(f"Going forward @ {speed}")
        await self.move_left(speed, 'forward')
        await self.move_right(speed, 'forward')

    async def go_backward(self, speed):
        logger.info(f"Going backward @ {speed}")
        await self.move_left(speed, 'backward')
        await self.move_right(speed, 'backward')

    async def turn_around(self, direction):
        logger.info(f"Turning around {direction}")
        if direction=='left':
            await self.move_left(-10, 'forward')
            await self.move_right(10, 'forward')
        else:
            await self.move_left(10, 'forward')
            await self.move_right(-10, 'forward')
