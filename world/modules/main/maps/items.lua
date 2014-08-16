item({
	id = "pot",
	name = "Potion of Awesome",
	type = "potion",
	variation = 0,
	description = "It is awesome",
	effect = function(origin, target)
		nudge(origin, origin, "health", -1)
		take(origin, "pot")
	end
})
weapon({
	id = "stabbysword",
	name = "Sword of stabbyness",
	type = "sword",
	variation = 0,
	description = "You could stab someone with this. Watch it!",
	effect = function(origin, target)
		nudge(origin, target, "health", -2)
	end
})
function createPotion(id, variation, name, amount)
	item({
		id = id,
		name = name,
		type = "potion",
		variation = variation,
		description = "Heals you some",
		effect = function(origin, target)
			nudge(origin, origin, "health", amount)
		end
	})
end
createPotion("hpotmin", 0, "Minor Health Potion", 2)
createPotion("hpotmed", 1, "Medium Health Potion", 5)
createPotion("hpotmax", 2, "Major Health Potion", 10)
item({
	id = "cream",
	name = "Amulet of Healing Potion Creation",
	type = "amulet",
	variation = 0,
	description = "Creates a Medium Health Potion",
	effect = function(origin, target)
		item("hpotmed", {
			map = "testmap",
			x = 6,
			y = 6,
		})
	end
})
item({
	id = "speaker",
	name = "Speaker",
	type = "speaker",
	variation = 0,
	description = "Say something rad",
	effect = function(origin, target)
		say("something rad", origin)
		say("something positionally rad", {map = "testmap", x = 2, y = 2})
	end
})
