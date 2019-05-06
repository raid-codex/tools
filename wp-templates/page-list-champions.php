<?php
/**
 * Template Name: Champion list
 *
 * The template for the Champion page.
 *
 */

require_once("tools/champion.php");
require_once("tools.php");

$data = champion_list();

get_header();

/**
 * Don't display page header if header layout is set as classic blog.
 */
do_action( 'hestia_before_single_page_wrapper' );

$championsPerFaction = array();

foreach ($data as $champion)
{
    $faction = $champion->{"faction"}->{"name"};
    if (!isset($championsPerFaction[$faction]))
    {
        $championsPerFaction[$faction] = array();
    }
    array_push($championsPerFaction[$faction], $champion);
}

$champions = array();
foreach ($championsPerFaction as $champions_in_faction)
{
    foreach ($champions_in_faction as $champion)
    {
        array_push($champions, $champion);
    }
}

?>

<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post champion-view">
		<div class="container">
            <article class="section pagebuilder-section">
                <div class="row champion-list-description align-left">
                    <div class="col-xs-12">
                        <p>
                            There are <?php echo sizeof($data); ?> champions indexed across <?php echo sizeof(array_keys($championsPerFaction)); ?> factions.
                        </p>
                    </div>
                </div>
                <div class="row">
                    <div class="col-xs-12">
                        <?php echo champion_get_html_table($champions, array("image","name","faction","rarity","type","element","rating_overall"), TRUE); ?>
                    </div>
                </div>
            </article>
        </div>
    </div>

	<?php get_footer(); ?>
