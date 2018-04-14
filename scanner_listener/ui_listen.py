#!/usr/bin/env python

import os
import sys
import time
import evdev
import requests
import json
import logging
from functools import partial
import fcntl

env = partial(os.environ.get)
log = logging.getLogger('barcode reader')
log.setLevel(logging.INFO)
formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")

ch = logging.StreamHandler(sys.stdout)
ch.setFormatter(formatter)
log.addHandler(ch)

class Config:
    SNACK_SERVICE_HOST_DEFAULT = "http://localhost:8080"
    SNACK_SERVICE_HOST = env("SNACK_SERVICE_HOST", SNACK_SERVICE_HOST_DEFAULT)
    INPUT_DEVICE_DIR_DEFAULT = "/dev/input"
    INPUT_DEVICE_DIR = env("INPUT_DEVICE_DIR", INPUT_DEVICE_DIR_DEFAULT)
    USER = env("USER")

class BarcodeReader():
    def __init__(self, config):
        log.info(config.INPUT_DEVICE_DIR)
        self.barcode_reader = None
        self.barcode_reader_addr = None
        self.config = config
        self.code = ''
        # Try to connect
        self.connect_to_reader()

    def post(self, code):
        log.info("Posting : {}".format(code))
        data = {"item_code": code}
        headers = {"content-type": "application/json"}
        log.info(self.config.SNACK_SERVICE_HOST )
        try:
            r = requests.post(self.config.SNACK_SERVICE_HOST  + "/state/item_code", data=json.dumps(data), headers=headers)
            print(r.content)
        except requests.exceptions.ConnectionError as e:
            print("Failed to connect. Please try again.")

    def find_barcode_reader(self):
        log.info("Listing devices.")
        devices = [evdev.InputDevice(fn) for fn in evdev.list_devices(self.config.INPUT_DEVICE_DIR)]

        for device in devices:
            log.info("{0}, {1}, {2}".format(device.fn, device.name, device.phys))
            if("Barcode Reader" in device.name):
                self.barcode_reader_addr = device.fn
                break

    def connect_to_reader(self):
        self.barcode_reader_addr = None
        self.barcode_reader = None
        try:
            self.find_barcode_reader()
            if self.barcode_reader_addr is not None:
                self.barcode_reader = evdev.InputDevice(self.barcode_reader_addr)
        except OSError as e:
            log.warning("No barcode reader found: {0}".format(e))

    def _read_code(self):
        for event in self.barcode_reader.read_loop():
            if event.type == evdev.ecodes.EV_KEY:
                key_press_data = evdev.categorize(event)
                # print(key_press_data)
                if key_press_data.keystate == 1:  # dow press only.
                    key_code = key_press_data.keycode[4:]
                    # print(key_press_data.keycode)
                    if len(key_code) == 1:
                        self.code = self.code + key_code
                    if key_code == "ENTER":
                        print("End of input")
                        print("Code is : {0}".format(self.code))
                        self.post(self.code)
                        self.code = ''

    def read_code(self):
        try:
            if(self.barcode_reader is not None):
                self._read_code()
            else:
                self.connect_to_reader()
                time.sleep(3)
        except IOError as e:
            log.warning("Missing input device.... is scanner plugged in? \n {0}".format(e))
            self.connect_to_reader()
            time.sleep(3)


def main():
    config = Config()
    reader = BarcodeReader(config)
    # devices = [evdev.InputDevice(fn) for fn in evdev.list_devices()]
    # for device in devices:
    #     print(device.fn, device.name, device.phys)
    while(1):
        reader.read_code()

if __name__ == "__main__":
    main()
