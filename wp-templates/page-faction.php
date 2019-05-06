<?php
/**
 * Template Name: Faction
 *
 * The template for the Faction page.
 *
 */

require_once("tools/champion.php");
require_once("tools/faction.php");
require_once("tools.php");

$custom = get_post_custom();
if (!isset($custom['faction-file']))
{
    wp_die("missing faction file");
}
$filename = $custom['faction-file'][0];
$faction = faction_get_by_filename($filename);

$champions_data = champion_list();

get_header();

/**
 * Don't display page header if header layout is set as classic blog.
 */
do_action( 'hestia_before_single_page_wrapper' );

$champions = array();
$factionSlug = explode(".", $filename)[0];
foreach ($champions_data as $champion)
{
    if ($champion->{"faction_slug"} == $factionSlug)
    {
        array_push($champions, $champion);
    }
}

?>

<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post champion-view">
		<div class="container">
            <article class="section pagebuilder-section">
                <div class="row faction-description align-left">
                    <div class="col-xs-12">
                        <p>
                            <?php echo $faction->{"name"}; ?> is a faction from RAID Shadow Legends composed of <?php echo sizeof($champions); ?> champions
                        </p>
                    </div>
                </div>
                <div class="row">
                    <div class="col-xs-12">
                        <h2>Champion list</h2>
                        <?php echo champion_get_html_table($champions, array("image","name","rarity","type","element","rating_overall")); ?>
                    </div>
                </div>
            </article>
        </div>
    </div>
    <?php get_footer(); ?>