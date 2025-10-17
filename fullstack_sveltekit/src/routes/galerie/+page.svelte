<script lang="ts">
    import {onMount} from "svelte";
    import {slide} from 'svelte/transition';
    import PicturesGallery from "$lib/PicturesGallery.svelte";
    import SideTitles from "$lib/SideTitles.svelte";

    /** @type {import('./$types').PageData} */
    export let data;

    let picturesContainer;

    let pictures;

    onMount(async () => {
        pictures = data.pictures;
    })
</script>

<main in:slide out:slide>

    <SideTitles/>

    {#if pictures}
        <div id="pictures" bind:this={picturesContainer}>
            <PicturesGallery pictures={pictures} mode="grouped"/>
        </div>
    {/if}
</main>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  *, *:before, *:after {
    box-sizing: border-box !important;
  }

  #pictures {
    position: absolute;
    top: 7vh;
    left: 10vw;
    width: 80vw;
  }
</style>
