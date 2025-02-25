<script lang="ts">
	import '../../app.css';
	let { children } = $props();
</script>

<nav class="top-menu">
	<div class="menu-title">Menu</div>
	<div class="hamburger-menu">
		<button>&#9776;</button>
	</div>
</nav>

<aside class="side-menu">
	<ul>
		<li><a href="#">Home</a></li>
		<li><a href="#">Playlist</a></li>
		<li><a href="#">Settings</a></li>
		<li><a href="#">About</a></li>
	</ul>
</aside>

{@render children()}
