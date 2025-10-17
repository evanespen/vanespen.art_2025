<script lang="ts">
    import {createEventDispatcher, onMount} from "svelte";
    import {scale} from 'svelte/transition';
    import {username} from "$lib/services/userStore.ts";
    import Icon from '@iconify/svelte';
    import {getHeaders} from "$lib/services/adminHeaders";

    export let picture;

    const dispatch = createEventDispatcher();
    let imageElem;
    let src = '';
    let classes = 'image-element';

    function isInViewport() {
        let rect = imageElem.getBoundingClientRect();
        if (
            rect.top >= 0 &&
            rect.left >= 0 &&
            rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) * 1.5 &&
            rect.right <= (window.innerWidth || document.documentElement.clientWidth) * 1.5
        ) {
            $: src = '/api/pictures/' + picture.path + '?type=thumb';
        }
    }


    function handleMouseMove(e) {
        imageElem.firstChild.style.transform = `scale(1.1) translateY(${e.clientY / 50 + 'px'}) translateX(${e.clientX / 50 + 'px'})`;
    }

    function handleMouseLeave(e) {
        imageElem.firstChild.style.transform = 'scale(1)';
    }

    function openLightbox() {
        dispatch('openLightbox', {picture});
    }

    async function starImage(evt) {
        evt.stopPropagation();
        const action = picture.stared ? 'unstar' : 'star';
        const response = await fetch(`/api/pictures?id=${picture.id}&action=${action}`, {headers: getHeaders()}, {
            method: 'PUT'
        });

        if (response.status === 200) {
            console.log('picture starred');
        }
    }

    async function deleteImage(evt) {
        evt.stopPropagation();
        const response = await fetch(`/api/pictures?id=${picture.id}`, {
            headers: getHeaders(),
            method: 'DELETE'
        });

        if (response.status === 200) {
            console.log('picture deleted')
        }
    }

    onMount(() => {
        if (picture.landscape) classes = classes + ' orient-landscape';
        else if (!picture.landscape) classes = classes + ' orient-portrait';
        else console.log('NO ORIENT');
        if (picture.stared) classes = classes + ' stared';
        if (picture.blured) classes = classes + ' blured';
        window.setTimeout(isInViewport, 300);
    });
</script>

<svelte:window on:scroll={isInViewport}/>

<div class={classes} transition:scale
     bind:this={imageElem}
     on:click={openLightbox}>
    <img {src} on:mousemove={handleMouseMove} on:mouseleave={handleMouseLeave}/>
    {#if $username}
        <div class="admin-buttons">
            <button class="btn-icon btn-warning" on:click={starImage}>
                <Icon icon="carbon:star" color="white" width={20}/>
            </button>
            <button class="btn-icon btn-error" on:click={deleteImage}>
                <Icon icon="carbon:trash-can" color="white" width={20}/>
            </button>
        </div>
    {/if}
</div>


<style lang="scss">
  .image-element {
    position: relative;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    min-height: 400px;
    //align-items: center;
    //justify-content: center;
    filter: drop-shadow(5px 5px 5px rgba(0, 0, 0, .3));

    .admin-buttons {
      position: absolute;
      bottom: 10px;
      right: 10px;
      z-index: 1000000;
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: calc(100% - 20px);

      button {
        margin-top: 10px;
        width: calc(50% - 5px);
      }
    }

    @media (max-width: 800px) {
      min-height: 200px;
    }

    &.orient-landscape {
      grid-column-end: span 2;
      grid-row-end: span 1;

      &.stared {
        @media (min-width: 800px) {
          grid-column-end: span 3;
          grid-row-end: span 2;
        }
      }
    }

    &.orient-portrait {
      grid-column-end: span 1;
      grid-row-end: span 1;

      &.stared {
        grid-column-end: span 2;
        grid-row-end: span 2;
      }
    }

    img {
      //height: 100%;
      //min-height: 400px;
      width: 100%;
      transition: .25s;
    }
  }

  //.blured {
  //  img {
  //    filter: blur(10px);
  //  }
  //}
</style>
