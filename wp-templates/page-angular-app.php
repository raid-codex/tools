<?php
/**
 * Template Name: Angular App
 *
 * The template for a page based on an Angular app
 *
 */

$custom = get_post_custom();
if (!isset($custom['angular-file-include']))
{
    wp_die("missing angular-file-include");
}
$ctrl = $custom['angular-file-include'][0];

get_header();

/**
 * Don't display page header if header layout is set as classic blog.
 */
do_action( 'hestia_before_single_page_wrapper' );

?>

<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post champion-view">
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
                        <div ng-app="websiteApp">
                            <div ng-include="'<?php echo $ctrl; ?>'"></div>
                        </div>
                    </div>
                </div>
            </article>
        </div>
    </div>
    <?php get_footer(); ?>