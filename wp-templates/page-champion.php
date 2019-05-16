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

wp_enqueue_style("champion-page", "https://raid-codex.com/wp-content/uploads/elementor/css/post-2747.css", FALSE, "1.0", "all");
wp_enqueue_style("elementor", "https://raid-codex.com/wp-content/plugins/elementor/assets/css/frontend.min.css?ver=2.5.15", FALSE, "2.5", "all");

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

?>

<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post champion-view">
		<div class="container">
        <article id="post-2747" class="section pagebuilder-section">
	<div data-elementor-type="post" data-elementor-id="2747" class="elementor elementor-2747" data-elementor-settings="[]">
		<div class="elementor-inner">
			<div class="elementor-section-wrap">
				<section class="elementor-element elementor-element-fb61089 elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-top-section" data-id="fb61089" data-element_type="section">
					<div class="elementor-container elementor-column-gap-default">
						<div class="elementor-row">
							<div class="elementor-element elementor-element-9e28071 elementor-column elementor-col-16 elementor-top-column" data-id="9e28071" data-element_type="column">
								<div class="elementor-column-wrap  elementor-element-populated">
									<div class="elementor-widget-wrap">
										<div class="elementor-element elementor-element-3c88ffa elementor-align-left elementor-widget elementor-widget-button" data-id="3c88ffa" data-element_type="widget" data-widget_type="button.default">
											<div class="elementor-widget-container">
												<div class="elementor-button-wrapper">
													<a href="https://raid-codex.com/champions/" class="elementor-button-link elementor-button elementor-size-xs" role="button">
													<span class="elementor-button-content-wrapper">
													<span class="elementor-button-text">
													<i class="fa fa-arrow-left"></i>
													</span>
													</span>
													</a>
												</div>
											</div>
										</div>
									</div>
								</div>
							</div>
							<div class="elementor-element elementor-element-3261f7a elementor-column elementor-col-66 elementor-top-column" data-id="3261f7a" data-element_type="column">
								<div class="elementor-column-wrap  elementor-element-populated">
									<div class="elementor-widget-wrap">
										<div class="elementor-element elementor-element-4f28b37 elementor-widget elementor-widget-heading" data-id="4f28b37" data-element_type="widget" data-widget_type="heading.default">
											<div class="elementor-widget-container">
											</div>
										</div>
									</div>
								</div>
							</div>
							<div class="elementor-element elementor-element-82f99c5 elementor-column elementor-col-16 elementor-top-column" data-id="82f99c5" data-element_type="column">
								<div class="elementor-column-wrap">
									<div class="elementor-widget-wrap">
									</div>
								</div>
							</div>
						</div>
					</div>
				</section>
				<section class="elementor-element elementor-element-cdf5d70 elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-top-section" data-id="cdf5d70" data-element_type="section">
					<div class="elementor-container elementor-column-gap-default">
						<div class="elementor-row">
							<div class="elementor-element elementor-element-258736c elementor-column elementor-col-100 elementor-top-column" data-id="258736c" data-element_type="column">
								<div class="elementor-column-wrap">
									<div class="elementor-widget-wrap">
									</div>
								</div>
							</div>
						</div>
					</div>
				</section>
				<section class="elementor-element elementor-element-3d88452 elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-top-section" data-id="3d88452" data-element_type="section">
					<div class="elementor-container elementor-column-gap-default">
						<div class="elementor-row">
							<div class="elementor-element elementor-element-ff34878 elementor-column elementor-col-50 elementor-top-column" data-id="ff34878" data-element_type="column">
								<div class="elementor-column-wrap  elementor-element-populated">
									<div class="elementor-widget-wrap">
										<div class="elementor-element elementor-element-f1b8f93 elementor-widget elementor-widget-image" data-id="f1b8f93" data-element_type="widget" data-widget_type="image.default">
											<div class="elementor-widget-container">
												<div class="elementor-image">
                                                    <?php echo get_image_url_by_slug($champion->{"image_slug"}, array(300, 300)); ?>
												</div>
											</div>
										</div>
									</div>
								</div>
							</div>
							<div class="elementor-element elementor-element-5698bad elementor-column elementor-col-50 elementor-top-column" data-id="5698bad" data-element_type="column">
								<div class="elementor-column-wrap  elementor-element-populated">
									<div class="elementor-widget-wrap">
										<section class="elementor-element elementor-element-3d3c627 elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-inner-section" data-id="3d3c627" data-element_type="section">
											<div class="elementor-container elementor-column-gap-default">
												<div class="elementor-row">
													<div class="elementor-element elementor-element-98ee396 elementor-column elementor-col-50 elementor-inner-column" data-id="98ee396" data-element_type="column">
														<div class="elementor-column-wrap  elementor-element-populated">
															<div class="elementor-widget-wrap">
																<div class="elementor-element elementor-element-6971856 elementor-widget elementor-widget-text-editor" data-id="6971856" data-element_type="widget" data-widget_type="text-editor.default">
																	<div class="elementor-widget-container">
																		<div class="elementor-text-editor elementor-clearfix">
																			<p style="text-align: center;"><span class="h2"><?php echo get_champion_rarity($champion->{"rarity"}); ?></span></p>
																		</div>
																	</div>
																</div>
															</div>
														</div>
													</div>
													<div class="elementor-element elementor-element-1eacc77 elementor-column elementor-col-50 elementor-inner-column" data-id="1eacc77" data-element_type="column">
														<div class="elementor-column-wrap  elementor-element-populated">
															<div class="elementor-widget-wrap">
																<div class="elementor-element elementor-element-8f3efde elementor-widget elementor-widget-text-editor" data-id="8f3efde" data-element_type="widget" data-widget_type="text-editor.default">
																	<div class="elementor-widget-container">
																		<div class="elementor-text-editor elementor-clearfix">
                                                                            <p style="text-align: center;">
                                                                            <span class="h2">
                                                                                <a href="<?php echo $champion->{"faction"}->{"website_link"}; ?>">
                                                                                    <?php echo $champion->{"faction"}->{"name"}; ?>
                                                                                </a>
</span>
                                                                            </p>
																		</div>
																	</div>
																</div>
															</div>
														</div>
													</div>
												</div>
											</div>
										</section>
										<div class="elementor-element elementor-element-d6e30b2 elementor-hidden-phone elementor-widget elementor-widget-spacer" data-id="d6e30b2" data-element_type="widget" data-widget_type="spacer.default">
											<div class="elementor-widget-container">
												<div class="elementor-spacer">
													<div class="elementor-spacer-inner"></div>
												</div>
											</div>
										</div>
										<div class="elementor-element elementor-element-34ee9dd elementor-widget elementor-widget-text-editor" data-id="34ee9dd" data-element_type="widget" data-widget_type="text-editor.default">
											<div class="elementor-widget-container">
												<div class="elementor-text-editor elementor-clearfix">
													<div class="centered">
                                                        <span class="h2">
                                                            <?php echo get_champion_grade($champion->{"rating"}->{"overall"}); ?>
                                                        </span>
													</div>
												</div>
											</div>
										</div>
										<section class="elementor-element elementor-element-bf0e946 elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-inner-section" data-id="bf0e946" data-element_type="section">
											<div class="elementor-container elementor-column-gap-default">
												<div class="elementor-row">
													<div class="elementor-element elementor-element-3aeed36 elementor-column elementor-col-50 elementor-inner-column" data-id="3aeed36" data-element_type="column">
														<div class="elementor-column-wrap  elementor-element-populated">
															<div class="elementor-widget-wrap">
																<div class="elementor-element elementor-element-a7bd9d3 elementor-widget elementor-widget-text-editor" data-id="a7bd9d3" data-element_type="widget" data-widget_type="text-editor.default">
																	<div class="elementor-widget-container">
																		<div class="elementor-text-editor elementor-clearfix">
																			<p style="text-align: center;"><strong>Element</strong>: <?php echo $champion->{"element"}; ?></p>
																		</div>
																	</div>
																</div>
															</div>
														</div>
													</div>
													<div class="elementor-element elementor-element-036d1f0 elementor-column elementor-col-50 elementor-inner-column" data-id="036d1f0" data-element_type="column">
														<div class="elementor-column-wrap  elementor-element-populated">
															<div class="elementor-widget-wrap">
																<div class="elementor-element elementor-element-e54ca88 elementor-widget elementor-widget-text-editor" data-id="e54ca88" data-element_type="widget" data-widget_type="text-editor.default">
																	<div class="elementor-widget-container">
																		<div class="elementor-text-editor elementor-clearfix">
																			<p style="text-align: center;"><b>Type</b>: <?php echo $champion->{"type"}; ?></p>
																		</div>
																	</div>
																</div>
															</div>
														</div>
													</div>
												</div>
											</div>
										</section>
									</div>
								</div>
							</div>
						</div>
					</div>
				</section>
				<section class="elementor-element elementor-element-d96ef5c elementor-reverse-mobile elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-top-section" data-id="d96ef5c" data-element_type="section">
					<div class="elementor-container elementor-column-gap-default">
						<div class="elementor-row">
							<div class="elementor-element elementor-element-729059d elementor-column elementor-col-50 elementor-top-column" data-id="729059d" data-element_type="column">
								<div class="elementor-column-wrap  elementor-element-populated">
									<div class="elementor-widget-wrap">
										<div class="elementor-element elementor-element-cd6574b elementor-widget elementor-widget-heading" data-id="cd6574b" data-element_type="widget" data-widget_type="heading.default">
											<div class="elementor-widget-container">
												<h3 class="elementor-heading-title elementor-size-default">Characteristics</h3>
											</div>
										</div>
										<div class="elementor-element elementor-element-7db2d0a elementor-tabs-view-horizontal elementor-widget elementor-widget-tabs" data-id="7db2d0a" data-element_type="widget" data-widget_type="tabs.default">
											<div class="elementor-widget-container">
												<div class="elementor-tabs" role="tablist">
													<div class="elementor-tabs-wrapper">
														<div id="elementor-tab-title-1311" class="elementor-tab-title elementor-tab-desktop-title elementor-active" data-tab="1" role="tab" aria-controls="elementor-tab-content-1311">
                                                            <a href="#characteristics-60">Level 60</a>
                                                        </div>
													</div>
													<div class="elementor-tabs-content-wrapper">
														<div class="elementor-tab-title elementor-tab-mobile-title elementor-active" data-tab="1" role="tab">Level 60</div>
														<div id="elementor-tab-content-1311" class="elementor-tab-content elementor-clearfix elementor-active" data-tab="1" role="tabpanel" aria-labelledby="elementor-tab-title-1311" style="display: block;">
															<table class="table table-hover table-no-border">
																<thead>
																	<tr>
																		<td colspan="2"></td>
																	</tr>
																</thead>
																<tbody>
                                                                <?php
                                                                $level = "60";
                                                                $champion_characteristics = $champion->{"characteristics"}->{$level};
                                                                $champion_characteristics_before = null;
                                                                if ($champion_pre17)
                                                                {
                                                                    $champion_characteristics_before = $champion_pre17->{"characteristics"}->{$level};
                                                                }
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
													</div>
												</div>
											</div>
										</div>
									</div>
								</div>
							</div>
							<div class="elementor-element elementor-element-4a81cf1 elementor-column elementor-col-50 elementor-top-column" data-id="4a81cf1" data-element_type="column">
								<div class="elementor-column-wrap  elementor-element-populated">
									<div class="elementor-widget-wrap">
										<div class="elementor-element elementor-element-9a9b121 elementor-widget elementor-widget-heading" data-id="9a9b121" data-element_type="widget" data-widget_type="heading.default">
											<div class="elementor-widget-container">
												<h3 class="elementor-heading-title elementor-size-default">Ratings</h3>
											</div>
										</div>
										<div class="elementor-element elementor-element-6a03b48 elementor-widget elementor-widget-text-editor" data-id="6a03b48" data-element_type="widget" data-widget_type="text-editor.default">
											<div class="elementor-widget-container">
												<div class="elementor-text-editor elementor-clearfix">
													<table class="table-hover table">
														<thead>
															<tr class="row-header">
																<th>Location</th>
																<th>Rating</th>
															</tr>
														</thead>
														<tbody>
                                                            <?php
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
                                                            ?>
														</tbody>
													</table>
												</div>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
                </section>
                <?php
                if (isset($champion->{"lore"}))
                {
                    ?>
                    <section class="elementor-element elementor-element-08c5503 elementor-section-boxed elementor-section-height-default elementor-section-height-default elementor-section elementor-top-section" data-id="08c5503" data-element_type="section">
                        <div class="elementor-container elementor-column-gap-default">
                            <div class="elementor-row">
                                <div class="elementor-element elementor-element-f3da6a5 elementor-column elementor-col-50 elementor-top-column" data-id="f3da6a5" data-element_type="column">
                                    <div class="elementor-column-wrap  elementor-element-populated">
                                        <div class="elementor-widget-wrap">
                                            <div class="elementor-element elementor-element-5bfde36 elementor-widget elementor-widget-heading" data-id="5bfde36" data-element_type="widget" data-widget_type="heading.default">
                                                <div class="elementor-widget-container">
                                                    <h2 class="elementor-heading-title elementor-size-default">Lore</h2>
                                                </div>
                                            </div>
                                            <div class="elementor-element elementor-element-14ee248 elementor-widget elementor-widget-text-editor" data-id="14ee248" data-element_type="widget" data-widget_type="text-editor.default">
                                                <div class="elementor-widget-container">
                                                    <div class="elementor-text-editor elementor-clearfix">
                                                        <p>
                                                            <?php echo $champion->{"lore"}; ?>
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="elementor-element elementor-element-e46daa0 elementor-column elementor-col-50 elementor-top-column" data-id="e46daa0" data-element_type="column">
                                    <div class="elementor-column-wrap">
                                        <div class="elementor-widget-wrap">
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </section>
                    <?php
                }
                ?>
			</div>
		</div>
	</div>
</article>
        </div>
    </div>

	<?php get_footer(); ?>
