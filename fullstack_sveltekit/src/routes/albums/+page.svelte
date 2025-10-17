<script>
    import {slide} from 'svelte/transition'
    import AlbumCard from '$lib/AlbumCard.svelte'
    import SideTitles from "$lib/SideTitles.svelte"

    /** @type {import('./$types').PageData} */
    export let data

</script>

<SideTitles/>

<main id="species" in:slide out:slide>
    {#each data.albums as album}
        <AlbumCard {album} kind="albums"/>
    {/each}
</main>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  #species {
    width: 80vw;
    margin-left: 10vw;
    margin-top: 10vh;

    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
  }
</style>
