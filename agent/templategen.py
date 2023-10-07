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

def make_gauge(gaugetype: str, entity_domain: str, entity_id: str, thresholds: list[float], title: str, pos: list[int]):
    return {
        "button": pos,
        "name": "homeassistant_string",
        "parameters": {
            "domain": entity_domain,
            "entity": entity_id,
            "icon": f"icons/gauge_{gaugetype}_full.png",
            "conditions": [
                {
                    "condition": {
                        "comparison": "<=",
                        "value": thresholds[0],
                    },
                    "icon": f"icons/gauge_{gaugetype}_empty.png",
                },
                {
                    "condition": {
                        "comparison": "<=",
                        "value": thresholds[1],
                    },
                    "icon": f"icons/gauge_{gaugetype}_low.png",
                },
                {
                    "condition": {
                        "comparison": "<=",
                        "value": thresholds[2],
                    },
                    "icon": f"icons/gauge_{gaugetype}_high.png",
                },
            ],
            "texts": [
                {
                    "color": [255, 255, 255, 255],
                    "size": 18,
                    "font": "font.ttf",
                    "x": 96/2,
                    "y": 0,
                    "align": "center",
                    "vertical-align": "top",
                    "text": title,
                },
                {
                    "color": [255, 255, 255, 255],
                    "size": 18,
                    "font": "font.ttf",
                    "x": 96/2,
                    "y": 96,
                    "align": "center",
                    "vertical-align": "bottom",
                    "text": "$STATE",
                }
            ],
        }
    }
    


PAGES = {}

def make_default_page():
    actions = []
    actions.append(make_fancy_light("light.hue_color_lamp_1_3", "floor_light_top", [0, 0]))
    actions.append(make_fancy_light("light.hue_color_candle_1_2", "floor_light_bottom", [0, 1]))
    actions.append(make_fancy_light("light.hue_lightguide_bulb_1", "ceiling_light", [1, 0]))

    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_co2", [600, 1000, 1500], "CO2", [2, 0]))
    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_particulate_matter_1_0um_concentration", [10, 50, 500], "1.0 um", [3, 0]))
    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_particulate_matter_2_5um_concentration", [10, 50, 500], "2.5 um", [4, 0]))
    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_particulate_matter_10_0um_concentration", [10, 50, 500], "10 um", [5, 0]))

    PAGES["default"] = {"actions":actions}

make_default_page()

with open("_gokrazy/extrafiles/etc/streamdeckpi/.gitignore", "w") as fign:
    for name, page in PAGES.items():
        fign.write(f"/{name}.yml\n")
        with open(f"_gokrazy/extrafiles/etc/streamdeckpi/{name}.yml", "w") as f:
            f.write(yaml_dump(page))
