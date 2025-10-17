<script lang="ts">
    import {createEventDispatcher, onMount, setContext} from "svelte";
    import Picture from "$lib/Picture.svelte";
    import LightBox from "$lib/LightBox.svelte";

    export let pictures;
    export let mode;
    let allPictures = [];
    let lightbox;

    const dispatch = createEventDispatcher();

    setContext('images-manager', {
        callPrevPicture: () => getPrevPicture,
        callNextPicture: () => getNextPicture,
    });

    function openLightbox(evt) {
        lightbox.setPicture(evt.detail.picture);
    }

    function getPrevPicture(current) {
        for (let i = 0; i < allPictures.length; i++) {
            if (allPictures[i].id === current.id) {
                if (i > 0) return allPictures[i - 1];
                else return current;
            }
        }
    }

    function getNextPicture(current) {
        for (let i = 0; i < allPictures.length; i++) {
            if (allPictures[i].id === current.id) {
                if (i < allPictures.length - 1) return allPictures[i + 1];
                else return current;
            }
        }
    }

    onMount(() => {
        if (mode == 'grouped') {
            Object.keys(pictures).forEach(month => {
                pictures[month].forEach(p => {
                    allPictures.push(p);
                })
            })
        } else if (mode == 'list') {
            allPictures = pictures;
        }
    })
</script>

<main>
    <LightBox bind:this={lightbox}/>

    {#if mode == 'grouped'}
        {#each Object.keys(pictures) as month}
            <div class="month-header">
                <h2>{month}</h2>
                <h5>({pictures[month].length} photos)</h5>
            </div>
            <div class="pictures-container">
                {#each pictures[month] as picture}
                    <Picture {picture} on:openLightbox={openLightbox}/>
                {/each}
            </div>
        {/each}
    {/if}

    {#if mode == 'list'}
        <div class="pictures-container">
            {#each pictures as picture}
                <Picture {picture} on:openLightbox={openLightbox}/>
            {/each}
        </div>
    {/if}
</main>


<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  main {
    margin-bottom: 10vh;
  }

  .month-header {
    border-bottom: 1px solid $text;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 3em;
    margin-bottom: 5vh;
    margin-top: 5vh;

    h2 {
      @include f-h-b;
      font-size: 2em;
    }

    h5 {
      @include f-p-b;
      font-style: oblique;
      font-size: 1.5em;
      color: $subcolor;
    }
  }

  .pictures-container {
    position: relative;
    display: grid;
    grid-template-columns: repeat(4, 18vw);
    //grid-auto-rows: 300px;
    //grid-template-columns: repeat(4, 20vw);
    //gap: 1vw;
    row-gap: 5vmin;
    column-gap: 5vmin;

    @media (max-width: 800px) {
      grid-template-columns: repeat(2, 37.5vw);
      //display: flex;
      //flex-direction: column;
    }

    img {
      height: 100%;
      min-height: 300px;
      width: 100%;
    }
  }
</style>
