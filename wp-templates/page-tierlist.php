<?php
/**
 * Template Name: Tier list
 *
 * The template for a Tier list page
 *
 */

require_once("tools/champion.php");
require_once("tools.php");

$custom = get_post_custom();
if (!isset($custom['tier-list']))
{
    wp_die("missing tier list");
}
$ratingKeys = explode(",", $custom['tier-list'][0]);
if (sizeof($ratingKeys) > 1 && !(substr($ratingKeys[0], 0, 6) == "arena_" || substr($ratingKeys[0], 0, 10) == "clan_boss_"))
{
    wp_die("invalid rating key");
}

$champions = champion_list();

$championsByKey = array();
foreach ($ratingKeys as $ratingKey)
{
    $championsByKey[$ratingKey] = array();
    foreach ($champions as $champion)
    {
        $key = $champion->{"rating"}->{$ratingKey};
        if (!isset($championsByKey[$ratingKey][$key]))
        {
            $championsByKey[$ratingKey][$key] = array();
        }
        array_push($championsByKey[$ratingKey][$key], $champion);
    }
}

get_header();

/**
 * Don't display page header if header layout is set as classic blog.
 */
do_action( 'hestia_before_single_page_wrapper' );

function get_section_display($rating, $title, $explanation_ok)
{
    global $championsByKey;
    global $ratingKeys;

    $hasSection = FALSE;
    foreach ($ratingKeys as $ratingKey)
    {
        $hasSection |= isset($championsByKey[$ratingKey][$rating]);
    }
    $section = <<<EOF
<div class="row">
    <div class="col-xs-12">
        <h2>$title</h2>
EOF;
    if (!$hasSection)
    {
        $section .= <<<EOF
            <p class="tier-explanation">None of the champions are ranked in this category</p>
EOF;
    }
    else
    {
        $section .= <<<EOF
        <p class="tier-explanation">
            $explanation_ok
        </p>
        <div class="row">
EOF;
        $class = sizeof($ratingKeys) == 1 ? "col-xs-12" : "col-xs-12 col-md-6";
        foreach ($ratingKeys as $ratingKey)
        {
            $table = champion_get_html_table($championsByKey[$ratingKey][$rating], array("image","name","faction","rarity","type","element","rating_".$ratingKey));
            $section .= <<<EOF
            <div class="$class">
                $table
            </div>
EOF;
        }
        $section .= <<<EOF
        </div>
EOF;
    }
    $section .= <<<EOF
    </div>
</div>
EOF;
    return $section;
}

/*function get_explanation_for_rating($rating)
{
    switch $rating {
        case "SS":
            return 
        case "S":
            return ""
        case "A":
            return "Those champions are considered as exceptional, their ability to deal with the fights is more than good but less than excellent."
        case "B":
            return "Those champions are considered as good, their ability to deal with the fights is good."
        case "C":
            return "Those champions are to avoid, their ability to deal with the fights is not that good."
        case "D":
            return "Those champions should not be brought at all in this fight."
    }
    return "";
}*/

?>

<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post tier-list-view">
		<div class="container">
            <article class="section pagebuilder-section">
                <div class="row">
                    <div class="col-xs-12">
                        <?php
                        wp_reset_query(); // necessary to reset query
                        while ( have_posts() ) : the_post();
                            the_content();
                        endwhile;
                        ?>
                    </div>
                </div>
                <div class="row">
                    <div class="col-xs-12">
                        <?php
                            echo get_section_display("SS", "God tier", "Those champions are considered as god tier, their ability to deal with the fights is above every other champion.");
                        ?>
                    </div>
                </div>
                <div class="row">
                    <div class="col-xs-12">
                        <?php
                            echo get_section_display("S", "Top tier", "Those champions are considered as top tier, their ability to deal with the fights is excellent.");
                        ?>
                    </div>
                </div>
            </article>
        </div>
    </div>
    <?php get_footer(); ?>