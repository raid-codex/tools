<?php
/**
 * Template Name: Champion
 *
 * The template for the Champion page.
 *
 */

require_once("tools/champion.php");
require_once("tools.php");

$custom = get_post_custom();
if (!isset($custom['champion-file']))
{
    die("missing champion file");
}
$filename = $custom['champion-file'][0];
$champion = champion_get_by_filename($filename);
$championSlug = $champion->{"slug"};
$description = champion_get_description($champion);

get_header();

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
                    <div class="col-xs-12" style="text-align: left;">
                        <a href="<?php echo get_permalink_by_slug("champions"); ?>">
                            <span class="btn btn-small"><i class="fa fa-arrow-left"></i></span>
                        </a>
                    </div>
                </div>
                <div class="row align-left champion-description">
                    <div class="col-xs-12">
                        <?php echo $description; ?>
                    </div>
                </div>
                <div class="row">
                    <div class="col-xs-12 col-md-6">
                        <div class="row">
                            <div class="col-xs-12">
                                <div class="row">
                                    <div class="col-xs-12 col-md-6">
                                        <?php echo get_image_url_by_slug($champion->{"image_slug"}); ?>
                                    </div>
                                    <div class="col-xs-12 col-md-6">
                                        <h2>
                                            <?php echo get_champion_grade($champion->{"rating"}->{"overall"}); ?>
                                        </h2>
                                    </div>
                                    <div class="col-xs-12 col-md-6">
                                        <h2>
                                            <?php echo get_champion_rarity($champion->{"rarity"}); ?>
                                        </h2>
                                    </div>
                                </div>
                                <div class="row">
                                    <div class="col-xs-12 col-lg-6 centered">
                                        <h4>
                                            <a href="<?php echo $champion->{"faction"}->{"website_link"}; ?>">
                                                <?php echo $champion->{"faction"}->{"name"}; ?>
                                            </a>
                                        </h4>
                                    </div>
                                    <div class="col-xs-12 col-lg-6 centered">
                                        <strong>Element:</strong> <?php echo $champion->{"element"}; ?><br>
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
                                    ?>
                                    <div id="characteristics-<?php echo $level; ?>">
                                        <table class="table table-hover table-responsive">
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
                                                            <?php
                                                            if ($champion_characteristics->{$c["key"]} == 0)
                                                            {
                                                                echo "Not specified";
                                                            }
                                                            else
                                                            {
                                                                switch ($c["type"]) {
                                                                    case "percentage":
                                                                        echo ($champion_characteristics->{$c["key"]} * 100)." %";
                                                                        break;
                                                                    default:
                                                                        echo $champion_characteristics->{$c["key"]};
                                                                        break;
                                                                }
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
                        <table class="table-hover table table-responsive">
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
                                        "display" => "Ice Guardian",
                                        "key" => "ice_guardian",
                                    ),
                                    array(
                                        "display" => "Dragon",
                                        "key" => "dragon",
                                    ),
                                    array(
                                        "display" => "Spider",
                                        "key" => "spider",
                                    ),
                                    array(
                                        "display" => "Fire Knight",
                                        "key" => "fire_knight",
                                    ),
                                    array(
                                        "display" => "Minotaur",
                                        "key" => "minotaur",
                                    ),
                                    array(
                                        "display" => "Force Dungeon",
                                        "key" => "force_dungeon",
                                    ),
                                    array(
                                        "display" => "Magic Dungeon",
                                        "key" => "magic_dungeon",
                                    ),
                                    array(
                                        "display" => "Spirit Dungeon",
                                        "key" => "spirit_dungeon",
                                    ),
                                    array(
                                        "display" => "Void Dungeon",
                                        "key" => "void_dungeon",
                                    )
                                );
                                foreach ($ratings as $rating_data) 
                                {
                                    $grade = $champion->{'rating'}->{$rating_data["key"]};
                                    ?>
                                    <tr>
                                        <td><?php echo $rating_data["display"]; ?></td>
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
            </article>
        </div>
    </div>

	<?php get_footer(); ?>
