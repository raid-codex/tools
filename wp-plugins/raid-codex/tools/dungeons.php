<?php

$keyToName = array(
    "ice_guardian" => "Ice Golem's Peak",
    "dragon" => "Dragon's Lair",
    "spider" => "Spider's Den",
    "fire_knight" => "Fire Knight's Castle",
    "minotaur" => "Minotaur's Labyrinth",
    "force_dungeon" => "Force Keep",
    "magic_dungeon" => "Magic Keep",
    "void_dungeon" => "Void Keep",
    "spirit_dungeon" => "Spirit Keep",
);

function dungeon_get_name_from_key($key)
{
    global $keyToName;
    if (isset($keyToName[$key]))
    {
        return $keyToName[$key];
    }
    return "";
}

?>