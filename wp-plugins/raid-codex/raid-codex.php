/**
* Plugin Name: Raid Codex
* Plugin URI: https://github.com/raid-codex/tools
* Description: Raid Codex tools
* Version: 1.0
* Author: Geoffrey Bauduin
* Author URI: https://github.com/geoffreybauduin
*/

add_action('save_post', 'save_post_hook', 10, 3);

function save_post_hook($post_ID, $post, $update) {
    if ($update)
    {
        return ; 
    }
	$post_type = get_post_type($post_ID);
  	if ($post_type == "page")
  	{
        $slug = $post->post_name;
        if (substr($slug, 0, 10) == "champions-")
        {   
            $exploded = explode("-", $slug);
            $filename = array_values(array_slice($exploded, -1))[0];
            add_post_meta( $post_ID, 'champion-file', $filename.'.json', true );
        }
    }
}

// if missing auth: RewriteRule .* - [E=REMOTE_USER:%{HTTP:Authorization}]
