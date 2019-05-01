<?php

function add_seo_metadata( $post )
{
    $description = "";
    update_post_meta( $post->id, '_yoast_wpseo_opengraph-title', '');
    update_post_meta( $post->id, '_yoast_wpseo_title', '%%title%% %%page%% %%sep%% %%parent_title%% %%sep%% %%sitename%%' );
    update_post_meta( $post->id, '_yoast_wpseo_focuskw', 'raid shadow legends  keyword2' );
    update_post_meta( $post->id, '_yoast_wpseo_metadesc', $description );
    update_post_meta( $post->id, '_yoast_wpseo_opengraph-description', $description );
}