import os
import playsound as ps


sound_base_path = "/mnt/c/Users/zero/Desktop/beater"
for f in os.listdir(sound_base_path):
    ps.playsound(os.path.join(sound_base_path, f))
