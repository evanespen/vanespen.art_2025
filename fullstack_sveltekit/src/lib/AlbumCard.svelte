<script lang="ts">
    export let album;
    export let kind;

    import {onMount} from "svelte";
    import {goto} from '$app/navigation';

    let specieCardImage;

    function goToAlbum() {
        goto(`/${kind}/${album.name}`);
    }

    onMount(() => {
        if (specieCardImage) {
            specieCardImage.style.backgroundImage = 'url(' + '/api/pictures/' + album.pictures[0].path + '?type=thumb' + ')';
        }
    })
</script>

{#if album}
    <main class="specie-card" on:click={goToAlbum}>
        <div class="specie-card-image" bind:this={specieCardImage}></div>
        <div class="specie-card-header">
            <div class="top">{album.name}</div>
            <div class="bottom">
                {#if kind == 'animalier'}
                    <div class="sci-name">{album.scientific_name}</div>
                {:else}
                    <div class="sci-name"></div>
                {/if}
                <div class="count">({album.pictures.length} photos)</div>
            </div>
        </div>
    </main>
{/if}

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  .specie-card {
    margin-bottom: 5vh;
    height: 30vh;
    width: 25vw;
    position: relative;
    overflow: hidden;

    @media (max-width: 800px) {
      width: 100%;
    }

    &:hover {
      cursor: pointer;

      .specie-card-image {
        transform: scale(1.1);
      }
    }

    .specie-card-image {
      position: relative;
      height: 100%;
      width: 100%;
      background-size: cover;
      background-position: center;
      z-index: -1;
      transition: .5s;
    }

    .specie-card-header {
      //@include f-h-b;
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      width: 100%;
      color: white;
      z-index: 10000;
      font-size: 1.7em;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      background-color: rgba(0, 0, 0, .3);

      .top {
        font-family: 'sterla' !important;
        font-variant: small-caps;
        padding: 10px;
        font-size: 2.5em;
        font-weight: bolder;
        text-align: left;
        filter: drop-shadow(2px 5px 5px rgba(0, 0, 0, .5));
      }

      .bottom {
        border-top: 2px solid white;
        padding: 10px;
        display: flex;
        justify-content: space-between;
        align-items: baseline;
        filter: drop-shadow(2px 5px 5px rgba(0, 0, 0, .5));

        .sci-name {
          font-style: oblique;
        }

        .count {
          @include f-p-b;
          font-style: oblique;
        }
      }
    }
  }
</style>
