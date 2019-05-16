<?php
/**
 * Template Name: Champion
 *
 * The template for the Champion page.
 *
 */

require_once("tools/champion.php");
require_once("tools.php");
require_once("tools/dungeons.php");

$custom = get_post_custom();
if (!isset($custom['champion-file']))
{
    wp_die("missing champion file");
}
$filename = $custom['champion-file'][0];
$champion = champion_get_by_filename($filename);
$championSlug = $champion->{"slug"};
$description = champion_get_description($champion);

$champion_pre17 = champion_get_by_filename($filename, "pre-1.7", FALSE);

get_header();

function get_characteristics_text($champion_characteristics, $c, $override_value=null)
{
    if ($champion_characteristics->{$c["key"]} == 0 && $champion_characteristics->{"hp"} == 0)
    {
        return "Not specified";
    }
    $value = $champion_characteristics->{$c["key"]};
    if ($override_value)
    {
        $value = $override_value;
    }
    if (!isset($c["type"])) { $c["type"] = "default"; }
    switch ($c["type"]) {
        case "percentage":
            return ($value * 100)." %";
    }
    return $value;
}

/**
 * Don't display page header if header layout is set as classic blog.
 */
do_action( 'hestia_before_single_page_wrapper' );

$characteristics = array(
    array(
        "display" => "HP",
        "key" => "hp",
    ),
    array(
        "display" => "Attack",
        "key" => "attack",
    ),
    array(
        "display" => "Defense",
        "key" => "defense",
    ),
    array(
        "display" => "Critical rate",
        "key" => "critical_rate",
        "type" => "percentage",
    ),
    array(
        "display" => "Critical damage",
        "key" => "critical_damage",
        "type" => "percentage",
    ),
    array(
        "display" => "Resistance",
        "key" => "resistance",
    ),
    array(
        "display" => "Accuracy",
        "key" => "accuracy",
    ),
);

?>

<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post champion-view">
		<div class="container">
            <article class="section pagebuilder-section centered">
                <div class="row">
                    <div class="col-xs-2 col-md-1" style="text-align: left;">
                        <a href="<?php echo get_permalink_by_slug("champions"); ?>">
                            <span class="btn btn-small"><i class="fa fa-arrow-left"></i></span>
                        </a>
                    </div>
                    <div class="col-xs-10 col-md-11">
                        <div class="row">
                            <div class="col-xs-12 col-md-6" style="margin-bottom: 15px;">
                                <span class="h2">
                                    <?php echo get_champion_rarity($champion->{"rarity"}); ?>
                                </span>
                            </div>
                            <div class="col-xs-12 col-md-6">
                                <span class="h2">
                                    <a href="<?php echo $champion->{"faction"}->{"website_link"}; ?>">
                                        <?php echo $champion->{"faction"}->{"name"}; ?>
                                    </a>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="row">
                    <div class="col-xs-12 col-md-6">
                        <div class="row">
                            <div class="col-xs-12">
                                <div class="row">
                                    <div class="col-xs-12 col-md-6">
                                        <?php echo get_image_url_by_slug($champion->{"image_slug"}, "medium"); ?>
                                    </div>
                                    <div class="col-xs-12 col-md-6">
                                        <span class="h2">
                                            <?php echo get_champion_grade($champion->{"rating"}->{"overall"}); ?>
                                        </span>
                                    </div>
                                </div>
                                <div class="row">
                                    <div class="col-xs-12 centered">
                                        <?php
                                        if ($champion->{"element"} != "")
                                        {
                                            ?>
                                            <strong>Element:</strong> <?php echo $champion->{"element"}; ?><br>
                                            <?php
                                        }
                                        ?>
                                        <strong>Type:</strong> <?php echo $champion->{"type"}; ?><br>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col-xs-12" id="characteristics-tabs">
                                <h2>Characteristics</h2>
                                <?php
                                foreach ($champion->{"characteristics"} as $level => $champion_characteristics)
                                {
                                    $champion_characteristics_before = null;
                                    if ($champion_pre17)
                                    {
                                        $champion_characteristics_before = $champion_pre17->{"characteristics"}->{$level};
                                    }
                                    ?>
                                    <div id="characteristics-<?php echo $level; ?>">
                                        <table class="table table-hover">
                                            <thead>
                                                <th colspan="2">
                                                    Level <?php echo $level; ?>
                                                </th>
                                            </thead>
                                            <tbody>
                                                <?php
                                                foreach ($characteristics as $c)
                                                {
                                                    ?>
                                                    <tr>
                                                        <td><?php echo $c["display"]; ?></td>
                                                        <td>
                                                            <span class="characteristic-current">
                                                                <?php
                                                                echo get_characteristics_text($champion_characteristics, $c);
                                                                ?>
                                                            </span>
                                                            <?php
                                                            if ($champion_characteristics_before && $champion_characteristics_before->{$c["key"]} != $champion_characteristics->{$c["key"]})
                                                            {
                                                                $type = ($champion_characteristics_before->{$c["key"]} > $champion_characteristics->{$c["key"]}) ? "down" : "up";
                                                                ?>
                                                                <span class="characteristic-before characteristic-before-<?php echo $type; ?>">
                                                                    (<i class="fas fa-arrow-<?php echo $type; ?> fa-rotate<?php echo ($type == "down" ? "-" : ""); ?>-45"></i> <?php echo get_characteristics_text($champion_characteristics_before, $c, abs($champion_characteristics_before->{$c["key"]} - $champion_characteristics->{$c["key"]})); ?> since 1.7)
                                                                </span>
                                                                <?php
                                                            }
                                                            ?>
                                                        </td>
                                                    </tr>
                                                    <?php
                                                    }
                                                    ?>
                                                </tbody>
                                            </table>
                                        </div>
                                        <?php
                                    }
                                ?>
                            </div>
                        </div>
                    </div>
                    <div class="col-xs-12 col-md-6">
                        <table class="table-hover table">
                            <thead>
                                <tr class="row-header">
                                    <th>Location</th>
                                    <th>Rating</th>
                                </tr>
                            </thead>
                            <tbody>
                                <?php
                                $ratings = array(
                                    array(
                                        "display" => "Campaign",
                                        "key" => "campaign",
                                    ),
                                    array(
                                        "display" => "Arena (off)",
                                        "key" => "arena_offense",
                                    ),
                                    array(
                                        "display" => "Arena (def)",
                                        "key" => "arena_defense",
                                    ),
                                    array(
                                        "display" => "Clan boss (without GS)",
                                        "key" => "clan_boss_without_giant_slayer",
                                    ),
                                    array(
                                        "display" => "Clan boss (with GS)",
                                        "key" => "clan_boss_with_giant_slayer",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "ice_guardian",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "dragon",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "spider",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "fire_knight",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "minotaur",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "force_dungeon",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "magic_dungeon",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "spirit_dungeon",
                                    ),
                                    array(
                                        "display" => null,
                                        "key" => "void_dungeon",
                                    )
                                );
                                foreach ($ratings as $rating_data) 
                                {
                                    $grade = $champion->{'rating'}->{$rating_data["key"]};
                                    ?>
                                    <tr>
                                        <td>
                                            <?php
                                            if ($rating_data["display"] == null)
                                            {
                                                $rating_data["display"] = dungeon_get_name_from_key($rating_data["key"]);
                                            }
                                            echo $rating_data["display"];
                                            ?>
                                        </td>
                                        <td class="champion-rating-<?php echo $grade; ?>"><?php echo getStarsForGrade($grade); ?></td>
                                    </tr>
                                    <?php
                                }
                                unset($rating_data);
                                ?>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="row align-left">
                    <div class="col-xs-12">
                        <?php
                        wp_reset_query(); // necessary to reset query
                        while ( have_posts() ) : the_post();
                            the_content();
                        endwhile;
                        ?>
                    </div>
                </div>
                <?php
                $page_content = get_the_content();
                if ($page_content == "")
                {
                    ?>
                    <div class="row align-left champion-description">
                        <div class="col-xs-12">
                            <?php echo $description; ?>
                        </div>
                    </div>
                    <?php
                }
                ?>
            </article>
        </div>
    </div>

	<?php get_footer(); ?>
