print("woah")

launcher.define_game{
  id = "test"
}

launcher.on_play(function(event)
  print(event.id, event.option)
  os.exec("echo launched game!!")
end)

print(launcher.name)