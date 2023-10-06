#!/usr/bin/env python3

from yaml import safe_dump as yaml_dump, safe_load as yaml_load

def make_fancy_light(entity_id: str, icon_type: str, pos: list[int]):
    return {
        "button": pos,
        "name": "multi",
        "parameters": {
            "render": {
                "name": "homeassistant_light",
                "parameters": {
                    "domain": "light",
                    "entity": entity_id,
                    "on_icon": f"icons/{icon_type}_on.png",
                    "off_icon": f"icons/{icon_type}_off.png",
                },
            },
            "run": [
                {
                    "name": "homeassistant_light",
                    "parameters": {
                        "domain": "light",
                        "entity": entity_id,
                        "service_name": "turn_on",
                    },
                    "pressed": True,
                    "conditions": [
                        {
                            "pressed": True,
                            "min": "10ms",
                            "max": "300ms",
                        },
                        {
                            "pressed": False,
                            "min": "100ms",
                        },
                    ],
                },
                {
                    "name": "homeassistant_light",
                    "parameters": {
                        "domain": "light",
                        "entity": entity_id,
                        "service_name": "turn_off",
                    },
                    "pressed": True,
                    "conditions": [
                        {
                            "pressed": True,
                            "min": "500ms",
                            "max": "1000ms",
                        },
                    ],
                },
            ],
        },
    }

PAGES = {}

def make_default_page():
    actions = []
    actions.append(make_fancy_light("light.hue_color_lamp_1_3", "floor_light_top", [0, 0]))
    actions.append(make_fancy_light("light.hue_color_candle_1_2", "floor_light_bottom", [0, 1]))
    actions.append(make_fancy_light("light.hue_lightguide_bulb_1", "ceiling_light", [1, 0]))
    PAGES["default"] = {"actions":actions}

make_default_page()

for name, page in PAGES.items():
    with open(f"_gokrazy/extrafiles/etc/streamdeckpi/{name}.yml", "w") as f:
        f.write(yaml_dump(page))
