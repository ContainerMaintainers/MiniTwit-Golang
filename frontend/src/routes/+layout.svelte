
<script lang="ts">
    export let cookie;
    console.log("cookie " + cookie)
    import '../app.css'
    import Navbar from "$lib/navbar.svelte"
</script>

<Navbar />
<slot />