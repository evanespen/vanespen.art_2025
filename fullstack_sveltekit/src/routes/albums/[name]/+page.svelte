<script>
    import PicturesGallery from "$lib/PicturesGallery.svelte";

    /** @type {import('./$types').PageData} */
    export let data;
    console.log('data', data);
</script>

{#if data.album}
    <main id="specie">
        <div id="specie-header"><h1>{data.album.name}</h1></div>

        <PicturesGallery pictures={data.album.pictures} mode="list"/>
    </main>
{/if}

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  #specie {
    width: 80vw;
    margin-top: 10vh;
    margin-left: 10vw;

    #specie-header {
      @include f-h-b;

      h1 {
        border-bottom: 1px solid black;
      }
    }
  }
</style>
