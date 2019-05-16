<?php
/**
 * Template Name: Champion - Generated
 *
 * The template for the Champion page.
 *
 */

wp_enqueue_style("champion-page", "https://raid-codex.com/wp-content/uploads/elementor/css/post-2747.css", FALSE, "1.0", "all");
wp_enqueue_style("elementor", "https://raid-codex.com/wp-content/plugins/elementor/assets/css/frontend.min.css?ver=2.5.15", FALSE, "2.5", "all");


get_header();

/**
 * Don't display page header if header layout is set as classic blog.
 */
do_action( 'hestia_before_single_page_wrapper' );

?>
<div class="<?php echo hestia_layout(); ?>">
	<div class="blog-post champion-view">
		<div class="container">
            <?php
            wp_reset_query(); // necessary to reset query
            while ( have_posts() ) : the_post();
                the_content();
            endwhile;
            ?>
        </div>
    </div>

	<?php get_footer(); ?>
