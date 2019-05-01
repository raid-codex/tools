<?php

require_once(__DIR__."/http.php");

$championUrl = "https://raid-codex.github.io/champions";

function champion_get_by_slug( $slug )
{
    return champion_get_by_filename("$slug.json");
}

function champion_get_by_filename ($filename)
{
    global $championUrl;

    $data = url_get("$championUrl/export/current/$filename");
    if ($data[1])
    {
        wp_die("cannot load champion: ".$data[1]->getMessage());
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
        $description = champion_get_default_description_as_text($champion, $lang);
        $content = "<p>$description</p>";
    }
    return $content;
}

function champion_get_default_description_as_text( $champion, $lang="en")
{
    return $champion->{"name"}." is a ".strtolower($champion->{"rarity"})." ".strtolower($champion->{"type"})." champion from the faction ".$champion->{"faction"}->{"name"}." doing ".strtolower($champion->{"element"})." damage";
}

?>