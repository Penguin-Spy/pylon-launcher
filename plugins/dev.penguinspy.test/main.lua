--[[ test.lua © Penguin_Spy 2024

  This Source Code Form is subject to the terms of the Mozilla Public
  License, v. 2.0. If a copy of the MPL was not distributed with this
  file, You can obtain one at http://mozilla.org/MPL/2.0/.
]]

print("woah this is lua")

launcher.define_game{
  id = "test",
  name = "Splatoon 2",
  hero = "hero.png"
}

launcher.define_game{
  id = "test2",
  name = "Among Us",
  hero = "hero2 - Copy.png"
}

launcher.on_play(function(event)
  print(event.id, event.option)
  os.execute("echo launched game!!")
end)

print("launcher.name = " .. tostring(launcher.name))

os.sleep(2)
print("finished loading")
