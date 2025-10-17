<script lang="ts">
    import {goto} from "$app/navigation";
    import {username} from "$lib/services/userStore.ts";
    import {page} from '$app/stores';


    let routes = {
        '/galerie': 'Galerie',
        '/animalier': 'Animalier',
        '/albums': 'Albums',
        // '/prestations': 'Prestations',
        // '/apropos': 'A Propos'
    };

    function logout() {
        goto('/logout');
    }
</script>

<main id="header">
    {#if $page.url.pathname != '/'}
        <a id="home-link-desktop" class="link" href="/">Accueil</a>
    {:else}
        <div id="home-link-desktop"></div>
    {/if}
    <div id="links">
        {#if $page.url.pathname != '/'}
            <a id="home-link-mobile" class="link" href="/">Accueil</a>
        {/if}
        {#each Object.keys(routes) as href}
            <a class:active={$page.url.pathname.includes(href)} class="link" {href}>{routes[href]}</a>
        {/each}
        {#if $username}
            <div class="logged-in" on:click={logout}>&times;</div>
        {/if}
    </div>
</main>

<style lang="scss">
  @import "$src/fonts.scss";
  @import "$src/color.scss";

  #header {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 5vh;
    font-size: 2em;
    border-bottom: 1px solid $text;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: $background;
    z-index: 100000;

    #home-link-desktop {
      margin-left: 1vw;
    }

    @media (max-width: 800px) {
      #home-link-desktop {
        display: none;
      }
    }

    #home-link-mobile {
      display: none;
    }

    @media (max-width: 800px) {
      #home-link-mobile {
        display: block;
      }
    }


    #links {
      width: 20%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-right: 50px;

      @media (max-width: 800px) {
        width: 100%;
        padding-right: 0;
        justify-content: space-around;
      }

      .logged-in {
        color: #F4511E;
        font-size: 1.5em;
        font-weight: bold;

        &:hover {
          cursor: pointer;
        }
      }
    }

    .link {
      @include f-p;
      transition: .5s;
      text-decoration: none;
      color: $text;

      &:hover {
        font-weight: 800;
        cursor: pointer;
      }
    }
  }

  :global(.active) {
    font-style: oblique !important;
    font-weight: bolder !important;
  }
</style>
