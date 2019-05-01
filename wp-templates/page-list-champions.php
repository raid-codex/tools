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
                        <table class="centered table-responsive table-hover table">
                            <thead>
                                <tr class="row-header">
                                    <th></th>
                                    <th>Champion name</th>
                                    <th>Faction</th>
                                    <th>Rarity</th>
                                    <th>Type</th>
                                    <th>Element</th>
                                    <th>Rank</th>
                                </tr>
                            </thead>
                            <tbody>
                                <?php
                                foreach ($championsPerFaction as $faction => $champions)
                                {
                                    ?>
                                    <tr class="row-header">
                                        <th colspan=7 class="centered">
                                            <?php echo $faction; ?>
                                        </th>
                                    </tr>
                                    <?php
                                    foreach ($champions as $champion_data)
                                    {
                                        ?>
                                        <tr>
                                            <td>
                                                <?php echo get_image_url_by_slug($champion_data->{"image_slug"}, array(30, 30)); ?>
                                            </td>
                                            <td>
                                                <strong>
                                                    <a href="<?php echo $champion_data->{"website_link"}; ?>">
                                                        <?php echo $champion_data->{"name"}; ?>
                                                    </a>
                                                </strong>
                                            </td>
                                            <td>
                                                <a href="<?php echo $champion_data->{"faction"}->{"website_link"}; ?>">
                                                    <?php echo $champion_data->{"faction"}->{"name"}; ?>
                                                </a>
                                            </td>
                                            <td>
                                                <?php echo get_champion_rarity($champion_data->{"rarity"}); ?>
                                            </td>
                                            <td>
                                                <?php echo $champion_data->{"type"}; ?>
                                            </td>
                                            <td>
                                                <?php echo $champion_data->{"element"}; ?>
                                            </td>
                                            <td>
                                                <?php echo get_champion_grade($champion_data->{"rating"}->{"overall"}); ?>
                                            </td>
                                        </tr>
                                        <?php
                                    }
                                }
                                ?>
                            </tbody>
                        </table>
                    </div>
                </div>
            </article>
        </div>
    </div>

	<?php get_footer(); ?>
