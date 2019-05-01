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
    die("missing faction file");
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
                        <table class="centered table-responsive table-hover table">
                            <thead>
                                <tr class="row-header">
                                    <th></th>
                                    <th>Champion name</th>
                                    <th>Rarity</th>
                                    <th>Type</th>
                                    <th>Element</th>
                                    <th>Rank</th>
                                </tr>
                            </thead>
                            <tbody>
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
                                ?>
                            </tbody>
                        </table>
                    </div>
                </div>
            </article>
        </div>
    </div>
</div>