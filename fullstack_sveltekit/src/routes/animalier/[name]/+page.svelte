<script>
    export let data

    import PicturesGallery from "$lib/PicturesGallery.svelte";
</script>

{#if data.specie}
    <main id="specie">
        <div id="specie-header"><h1>{data.specie.name}</h1></div>

        <PicturesGallery pictures={data.specie.pictures} mode="list" />
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
