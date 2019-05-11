<?php

require_once(__DIR__."/http.php");
require_once(__DIR__."/dungeons.php");

$championUrl = "https://raid-codex.github.io/champions";

function champion_get_by_slug( $slug )
{
    return champion_get_by_filename("$slug.json");
}

function champion_get_by_filename ($filename, $from="current", $die_if_error=TRUE)
{
    global $championUrl;

    $data = url_get("$championUrl/export/$from/$filename");
    if ($data[1])
    {
        if ($die_if_error)
        {
            wp_die("cannot load champion: ".$data[1]->getMessage());
        }
        return null;
    }
    $champion = json_decode($data[0]);
    return $champion;
}

function champion_list()
{
    return champion_get_by_filename("index.json");
}

function champion_get_description ( $champion, $lang="en" )
{
    global $championUrl;

    $content = url_get("$championUrl/descriptions/current/en/".$champion->{"slug"}.".html")[0];
    if (!$content)
    {
        $description = $champion->{"default_description"};
        $content = "<p>$description</p>";
    }
    return $content;
}

function champion_get_default_description_as_text( $champion, $lang="en")
{
    return $champion->{"name"}." is a ".strtolower($champion->{"rarity"})." ".strtolower($champion->{"type"})." champion from the faction ".$champion->{"faction"}->{"name"}." doing ".strtolower($champion->{"element"})." damage";
}

$champions_field_array = array(
    "image" => array(
        "display" => function ($champion) { return get_image_url_by_slug($champion->{"image_slug"}, "thumbnail"); },
        "header_name" => "",
    ),
    "name" => array(
        "display" => function ($champion) {
            return '<strong><a href="'.$champion->{"website_link"}.'">'.$champion->{"name"}.'</a></strong>';
        },
        "header_name" => "Name",
    ),
    "faction" => array(
        "display" => function ($champion) {
            return '<a href="'.$champion->{"faction"}->{"website_link"}.'">'.$champion->{"faction"}->{"name"}.'</a>';
        },
        "header_name" => "Faction",
    ),
    "rarity" => array(
        "display" => function ($champion) {
            return get_champion_rarity($champion->{"rarity"});
        },
        "header_name" => "Rarity",
    ),
    "type" => array(
        "display" => function ($champion) { return $champion->{"type"}; },
        "header_name" => "Type",
    ),
    "element" => array(
        "display" => function ($champion) { return $champion->{"element"}; },
        "header_name" => "Element",
    ),
    "rating_overall" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"overall"}); },
        "header_name" => "Overall rating",
    ),
    "rating_campaign" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"campaign"}); },
        "header_name" => "Campaign rating",
    ),
    "rating_arena_offense" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"arena_offense"}); },
        "header_name" => "Arena (offensive) rating",
    ),
    "rating_arena_defense" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"arena_defense"}); },
        "header_name" => "Arena (defense) rating",
    ),
    "rating_clan_boss_without_giant_slayer" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"clan_boss_without_giant_slayer"}); },
        "header_name" => "Clan Boss (without Giant Slayer) rating",
    ),
    "rating_clan_boss_with_giant_slayer" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"clan_boss_with_giant_slayer"}); },
        "header_name" => "Clan Boss (with Giant Slayer) rating",
    ),
    "rating_ice_guardian" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"ice_guardian"}); },
        "header_name" => dungeon_get_name_from_key("ice_guardian")." rating",
    ),
    "rating_dragon" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"dragon"}); },
        "header_name" => dungeon_get_name_from_key("dragon")." rating",
    ),
    "rating_spider" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"spider"}); },
        "header_name" => dungeon_get_name_from_key("spider")." rating",
    ),
    "rating_fire_knight" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"fire_knight"}); },
        "header_name" => dungeon_get_name_from_key("fire_knight")." rating",
    ),
    "rating_minotaur" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"minotaur"}); },
        "header_name" => dungeon_get_name_from_key("minotaur")." rating",
    ),
    "rating_force_dungeon" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"force_dungeon"}); },
        "header_name" => dungeon_get_name_from_key("force_dungeon")." rating",
    ),
    "rating_magic_dungeon" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"magic_dungeon"}); },
        "header_name" => dungeon_get_name_from_key("magic_dungeon")." rating",
    ),
    "rating_spirit_dungeon" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"spirit_dungeon"}); },
        "header_name" => dungeon_get_name_from_key("spirit_dungeon")." rating",
    ),
    "rating_void_dungeon" => array(
        "display" => function ($champion) { return get_champion_grade($champion->{"rating"}->{"void_dungeon"}); },
        "header_name" => dungeon_get_name_from_key("void_dungeon")." rating",
    ),
);

function champion_get_html_table($champions, $fields, $separate_factions=FALSE)
{
    global $champions_field_array;

    $lastFaction = null;
    $phoneTable = '<table class="centered table-responsive table-hover table champion-list-table no-header-mobile hidden-sm hidden-md hidden-lg">';
    $table = '<table class="centered table-responsive table-hover table champion-list-table no-header-mobile hidden-xs"><thead><tr class="row-header">';
    foreach ($fields as $field)
    {
        $header = $champions_field_array[$field]["header_name"];
        $table .= '<th class="table-header-'.$field.'">'.$header.'</th>';
    }
    $table .= '</tr></thead><tbody>';
    $phoneTable .= '<tbody>';
    foreach ($champions as $champion)
    {
        if ($champion->{"faction_slug"} != $lastFaction && $separate_factions)
        {
            $table .= '<tr class="row-header"><th colspan='.sizeof($fields).' class="centered">'.$champion->{"faction"}->{"name"}.'</th></tr>';
            $phoneTable .= '<tr class="row-header"><th class="centered">'.$champion->{"faction"}->{"name"}.'</th></tr>';
        }
        $lastFaction = $champion->{"faction_slug"};
        $table .= '<tr>';
        $phoneCell = "";
        foreach ($fields as $field)
        {
            $cellValue = $champions_field_array[$field]["display"]($champion);
            $table .= '<td class="table-col-'.$field.'">'.$cellValue.'</td>';
            $phoneCell .= "$cellValue<br/>";
        }
        $table .= "</tr>";
        $phoneTable .= '<tr><td>'.$phoneCell.'</td></tr>';
    }
    $table .= "</tbody></table>$phoneTable</tbody></table>";
    return $table;
}

?>