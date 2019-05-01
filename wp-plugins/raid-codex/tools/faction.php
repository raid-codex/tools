<?php

require_once(__DIR__."/http.php");

$factionUrl = "https://raid-codex.github.io/factions";

function faction_get_by_slug( $slug )
{
    return faction_get_by_filename("$slug.json");
}

function faction_get_by_filename ($filename)
{
    global $factionUrl;

    $data = url_get("$factionUrl/export/current/$filename");
    if ($data[1])
    {
        wp_die("cannot load faction: ".$data[1]->getMessage());
    }
    $faction = json_decode($data[0]);
    return $faction;
}

function faction_list()
{
    return faction_get_by_filename("index.json");
}

?>