<script lang="ts">
    import {beforeUpdate, createEventDispatcher, getContext} from "svelte";
    import Moment from 'moment';
    import Close from '$lib/svgs/close.svg';
    import ChevronLeft from '$lib/svgs/chevronLeft.svg';
    import ChevronRight from '$lib/svgs/chevronRight.svg';
    import Aperture from '$lib/svgs/aperture.svg?url';
    import Calendar from '$lib/svgs/calendar.svg?url';
    import Camera from '$lib/svgs/camera.svg?url';
    import Flash from '$lib/svgs/flash.svg?url';
    import Iso from '$lib/svgs/iso.svg?url';
    import Ruler from '$lib/svgs/ruler.svg?url';
    import Settings from '$lib/svgs/settings.svg?url';
    import Timer from '$lib/svgs/time.svg?url';
    import Chemistry from '$lib/svgs/chemistry.svg?url';
    import Note from '$lib/svgs/note.svg?url';
    import Star from '$lib/svgs/star.svg?url';

    const dispatch = createEventDispatcher()
    let lightBox
    let picture
    let exif

    const modes = {
        'Aperture priority': 'Priorité ouverture',
        'Shutter priority': 'Priorité vitesse',
        'Manual': 'Manuel',
        'Auto': 'Auto'
    }

    const {callPrevPicture, callNextPicture} = getContext('images-manager')

    Moment.locale('fr', {
        months: ['Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Août', 'Septembre', 'Octobre', 'Novembre', 'Décembre']
    })

    export function setPicture(_picture) {
        picture = _picture
        console.log(picture)


        if (picture.notes !== '' && picture.notes !== null) {
            if (picture.notes.split('XXX')[0] === 'ASTRO') {
                exif = [
                    {icon: Calendar, value: Moment(picture.day).format('DD MMMM YYYY', 'fr')},
                    {icon: Star, value: picture.notes.split('XXX')[1]},
                    {icon: Note, value: picture.notes.split('XXX')[2]},
                ];
            } else if (picture.notes.split('XXX')[0] === 'FILM') {
                exif = [
                    {icon: Calendar, value: Moment(picture.day).format('DD MMMM YYYY', 'fr')},
                    {icon: Chemistry, value: picture.notes.split('XXX')[1]},
                    {
                        icon: Note,
                        value: picture.notes.split('XXX')[2] !== 'null' ? picture.notes.split('XXX')[2] : ''
                    },
                ];
            }

        } else {
            console.log('normal picture')
            exif = [
                {icon: Calendar, value: Moment(picture?.day).format('DD MMMM YYYY', 'fr')},
                {icon: Camera, value: picture?.cam_model},
                {icon: Settings, value: modes[picture?.mode]},
                {icon: Aperture, value: picture?.aperture},
                {icon: Iso, value: picture?.iso},
                {icon: Timer, value: picture?.exposure},
                {icon: Ruler, value: picture?.focal},
                {icon: Flash, value: picture?.flash.includes('No') ? 'Non' : 'Oui'},
            ]
        }
    }

    export function closeLightbox() {
        if (lightBox == null) return
        lightBox.classList.remove('opening')
        lightBox.classList.add('closing')
        lightBox.style.left = '-100vw'
        window.setTimeout(() => {
            picture = undefined
            exif = undefined
            dispatch('closeLightbox', {})
            lightBox.style.left = '0'
            lightBox.classList.remove('closing')
            lightBox.classList.add('opening')
        }, 250)
    }

    function prevPicture() {
        picture ? setPicture(callPrevPicture()(picture)) : ''
    }

    function nextPicture() {
        picture ? setPicture(callNextPicture()(picture)) : ''
    }

    function handleKeyNavigation(event) {
        switch (event.key) {
            case 'ArrowLeft':
                prevPicture()
                break

            case 'ArrowRight':
                nextPicture()
                break

            case 'Escape':
                closeLightbox()
                break

            default:
                break
        }
    }

    beforeUpdate(() => {
        if (window && lightBox) {
            import('hammerjs').then(Hammer => {
                const manager = new Hammer.Manager(lightBox);
                const Swipe = new Hammer.Swipe();
                manager.add(Swipe);
                manager.on('swipe', function (e) {
                    if (e.offsetDirection === 2) {
                        nextPicture();
                    } else if (e.offsetDirection === 4) {
                        prevPicture();
                    }
                });
            });
        }
    });
</script>


<svelte:window on:keydown={handleKeyNavigation}/>

{#if picture}
    <main bind:this={lightBox} class="opening">
        <button id="close-btn" on:click={closeLightbox}>
            <Close height="50px"/>
        </button>

        <button id="prev-btn" on:click={prevPicture}>
            <ChevronLeft height="50px"/>
        </button>
        <button id="next-btn" on:click={nextPicture}>
            <ChevronRight height="50px"/>
        </button>

        <div id="lightbox-content">

            <div id="lightbox-content-image-container">
                <img src={'/api/pictures/' + picture.path + '?type=half'} alt="">
            </div>

            <div id="lightbox-content-infos">
                {#each exif as exifEntry}
                    <div class="lightbox-infos-entry">
                        <img class="icon" src={exifEntry.icon} alt="">
                        <div class="value">{exifEntry.value}</div>
                    </div>
                {/each}
            </div>
        </div>
    </main>
{/if}

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  @keyframes slide-from-left {
    0% {
      left: -100vw;
    }
    100% {
      left: 0;
    }
  }

  @keyframes slide-to-left {
    0% {
      left: 0;
    }
    100% {
      left: -100vw !important;
    }
  }

  main {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, .8);
    z-index: 1000000;
    backdrop-filter: blur(10px);

    ._button {
      background: none;
      border: none;
      outline: none;
      filter: invert(1);
      position: fixed;
      transition: .5s;

      &:hover {
        cursor: pointer;
      }
    }

    #close-btn {
      @extend ._button;
      top: 2vmin;
      right: 3vmin;

      &:hover {
        transform: scale(1.5) rotate(90deg);
      }

    }

    @keyframes slide {
      0% {
        transform: translateX(0);
      }
      50% {
        transform: translateX(10px);
      }
    }

    #prev-btn {
      @extend ._button;
      bottom: 2vmin;
      left: 2vmin;

      &:hover {
        animation: slide 1s infinite;
      }
    }

    #next-btn {
      @extend ._button;
      bottom: 2vmin;
      right: 3vmin;

      &:hover {
        animation: slide 1s infinite;
      }
    }

    #lightbox-content {
      position: fixed;
      top: 5vh;
      left: 2.5vw;
      height: 90vh;
      width: 95vw;

      filter: drop-shadow(2px 5px 5px rgba(0, 0, 0, .9));
      display: flex;
      flex-direction: row;
      align-items: flex-end;
      justify-content: space-between;

      #lightbox-content-infos {
        display: flex;
        flex-direction: column;
        width: 15vw;

        @media (max-width: 800px) {
          display: none;
        }

        .lightbox-infos-entry {
          display: flex;
          flex-direction: row;

          .icon {
            height: 30px;
            filter: invert(1);
            margin-right: 30px;
            margin-left: 30px;
          }

          .value {
            @include f-p-2;
            font-size: 1.5em;
            color: white;
            filter: drop-shadow(2px 5px 5px rgba(0, 0, 0, .9));
          }
        }
      }

      #lightbox-content-image-container {
        width: 80vw;
        height: 90vh;
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;

        @media (max-width: 800px) {
          width: 100vw;
        }

        img {
          max-width: 75vw;
          max-height: 90vh;
        }
      }
    }
  }

  :global(.opening) {
    animation: slide-from-left .25s !important;
    animation-play-state: running !important;
  }

  :global(.closing) {
    animation: slide-to-left .25s !important;
    animation-play-state: running !important;
  }
</style>
