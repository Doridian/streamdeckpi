#!/usr/bin/env python3

from yaml import safe_dump as yaml_dump, SafeDumper

SafeDumper.ignore_aliases = lambda *args : True

def make_onoff(entity_domain: str, entity_id: str, icon_type: str, action_type: str, pos: list[int]):
    return {
        "button": pos,
        "name": "multi",
        "parameters": {
            "render": {
                "name": action_type,
                "parameters": {
                    "domain": entity_domain,
                    "entity": entity_id,
                },
            },
            "run": [
                {
                    "name": action_type,
                    "parameters": {
                        "domain": entity_domain,
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
                    "name": action_type,
                    "parameters": {
                        "domain": entity_domain,
                        "entity": entity_id,
                        "service_name": "turn_off",
                    },
                    "pressed": True,
                    "conditions": [
                        {
                            "pressed": True,
                            "min": "500ms",
                            "max": "2000ms",
                        },
                    ],
                },
            ],
        },
    }

def make_light(entity_id: str, icon_type: str, pos: list[int]):
    action = make_onoff("light", entity_id, icon_type, "homeassistant_light", pos)
    params = action["parameters"]["render"]["parameters"]
    params["on_icon"] = f"icons/{icon_type}_on.png"
    params["off_icon"] = f"icons/{icon_type}_off.png"
    return action

def make_switch(entity_id: str, icon_type: str, pos: list[int]):
    action = make_onoff("switch", entity_id, icon_type, "homeassistant_entity", pos)
    params = action["parameters"]["render"]["parameters"]
    params["icon"] = f"icons/{icon_type}_off.png"
    params["conditions"] = [
        {
            "condition": {
                "comparison": "==",
                "value": "on",
            },
            "icon": f"icons/{icon_type}_on.png",
        }
    ]
    return action

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

def make_light_subpage(entity_id: str, icon_type: str, pos: list[int]):
    subpage_name = f"light_{entity_id}"

    actions = []
    actions.append(make_light(entity_id, icon_type, [0,0]))
    actions.append({
        "button": [7, 0],
        "name": "pop_page",
        "parameters": {
            "icon": "icons/back.png",
        }
    })
    for i in range(0, 8):
        brightness = i*32
        actions.append({
            "button": [i, 1],
            "name": "homeassistant_light",
            "parameters": {
                "on_icon": f"icons/{icon_type}_on.png",
                "off_icon": f"icons/{icon_type}_off.png",
                "domain": "light",
                "entity": entity_id,
                "service_name": "turn_on" if i > 0 else "turn_off",
                "service_data": {
                    "brightness": brightness,
                } if i > 0 else None,
                "render_state": "on" if i > 0 else "off",
                "render_brightness": brightness,
            }
        })

    preset_colors = [
        [255, 0  , 0  ],
        [0  , 255, 0  ],
        [0  , 0  , 255],
        [255, 255, 0  ],
        [0  , 255, 255],
        [255, 0  , 255],
        [255, 128, 0  ],
        [255, 255, 255],

        [255, 147, 41 ],
        [201, 226 ,255],
        [64 , 156, 255],
        [255, 244, 229],
        [244, 255, 250],
        [255, 197, 143],
        [255, 214, 170],
        [255, 241, 224],
    ]
    for i in range(0, 16):
        x = i % 16
        y = i // 16

        rgb_color = preset_colors[i]

        actions.append({
            "button": [x, y + 2],
            "name": "homeassistant_light",
            "parameters": {
                "on_icon": f"icons/{icon_type}_on.png",
                "off_icon": f"icons/{icon_type}_off.png",
                "domain": "light",
                "entity": entity_id,
                "service_name": "turn_on",
                "service_data": {
                    "rgb_color": rgb_color,
                },
                "render_state": "on",
                "render_rgb_color": rgb_color,
            }
        })

    PAGES[subpage_name] = {"actions":actions}

    action = make_light(entity_id, icon_type, pos)
    long_action = action["parameters"]["run"][1]
    long_action["name"] = "push_page"
    long_action["parameters"] = {
        "target": f"{subpage_name}.yml",
    }
    return action

def make_default_page():
    actions = []
    actions.append(make_light_subpage("light.hue_color_lamp_1_3", "floor_light_top", [1, 0]))
    actions.append(make_light_subpage("light.hue_color_candle_1_2", "floor_light_bottom", [1, 1]))
    actions.append(make_light_subpage("light.hue_lightguide_bulb_1", "ceiling_light", [2, 0]))
    actions.append(make_light_subpage("light.shapes_ea58", "triangle", [2, 1]))

    actions.append(make_light_subpage("light.dori_office", "light_group", [1, 2]))

    actions.append(make_switch("switch.dori_pc_switch", "desktop", [1, 3]))
    actions.append(make_switch("switch.dori_desktop_relay", "monitor", [2, 3]))

    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_co2", [600, 1000, 1500], "CO2", [0, 0]))
    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_particulate_matter_1_0um_concentration", [10, 50, 500], "1.0 um", [0, 1]))
    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_particulate_matter_2_5um_concentration", [10, 50, 500], "2.5 um", [0, 2]))
    actions.append(make_gauge("lowcolor", "sensor", "sensor.dori_office_particulate_matter_10_0um_concentration", [10, 50, 500], "10 um", [0, 3]))

    PAGES["default"] = {"actions":actions}

make_default_page()

with open("_gokrazy/extrafiles/etc/streamdeckpi/.gitignore", "w") as fign:
    for name, page in PAGES.items():
        fign.write(f"/{name}.yml\n")
        with open(f"_gokrazy/extrafiles/etc/streamdeckpi/{name}.yml", "w") as f:
            f.write(yaml_dump(page))
