<script lang="ts">
    import {onMount} from "svelte"
    import Renew from '$lib/svgs/renew.svg?url';
    import Trashcan from '$lib/svgs/trashcan.svg?url';
    import View from '$lib/svgs/view.svg?url';
    import List from '$lib/svgs/list.svg?url';
    import {goto} from "$app/navigation";
    import Status from '$lib/Status.svelte';
    import {getHeaders} from "$lib/services/adminHeaders";

    let newReview = {
        name: '',
        password: ''
    }
    let reviews = [];
    let showEvents = false;
    let events = {}, archivesDone = undefined;

    function processEvents(rawEvents): boolean {
        let _events = events;
        let done = false;
        for (const event of rawEvents) {
            if (event.fileName === 'archive') {
                archivesDone = event.status;
            } else if (event.fileName === 'ALL' && event.status === true) {
                done = true;
            } else {
                _events[event.fileName][event.step] = event.status;
            }
        }

        $: events = _events;
        return done;
    }

    function loadReviews() {
        fetch('/api/reviews', {headers: getHeaders()}).then(res => {
            res.json().then(data => {
                $: reviews = data.reviews
            })
        })
    }

    function handleNewReview() {
        fetch('/api/reviews', {
            method: 'POST',
            body: JSON.stringify(newReview),
            headers: getHeaders(),
        }).then(res => {
            console.log(res.status);
            loadReviews()
        })
    }

    async function handleLoadEvents(name) {
        fetch(`/api/reviews/${name}`, {
            method: 'GET',
            headers: getHeaders(),
        }).then(async res => {
            const body = await res.json();
            let _events = {};
            showEvents = true;
            body.review.pictures.forEach(picture => {
                _events[picture.name] = {
                    half: undefined,
                    db: undefined,
                };
                $: events = _events;
            });

            let done = false;
            let pollingLoop = window.setInterval(async () => {
                await fetch(`/api/reviews/${name}/events`, {headers: getHeaders()}).then(async res => {
                    const rawEvents = await res.json();
                    if (rawEvents.events.length > 0) {
                        done = processEvents(rawEvents.events);
                        if (done) clearInterval(pollingLoop);
                    }
                })
            }, 1000);
        });
    }

    async function handleRefresh(name) {
        showEvents = true;
        fetch(`/api/reviews/${name}`, {
            method: 'PUT',
            headers: getHeaders(),
            body: JSON.stringify({
                action: 'refresh'
            })
        }).then(async res => {
            const body = await res.json();
            const picturesList = body.picturesList;
            let _events = {};
            picturesList.forEach(pictureName => {
                _events[pictureName] = {
                    half: undefined,
                    db: undefined,
                };
                $: events = _events;
            });

            let done = false;
            let pollingLoop = window.setInterval(async () => {
                await fetch(`/api/reviews/${name}/events`, {headers: getHeaders()}).then(async res => {
                    const rawEvents = await res.json();
                    if (rawEvents.events.length > 0) {
                        done = processEvents(rawEvents.events);
                        if (done) clearInterval(pollingLoop);
                    }
                })
            }, 1000);
        })
    }

    function handleDelete(name) {
        fetch(`/api/reviews/${name}`, {
            method: 'DELETE',
            headers: getHeaders(),
        }).then(res => {
            console.log(res.status);
            loadReviews();
        })
    }

    function handleView(name) {
        goto(`/reviews/${name}`)
    }

    onMount(() => {
        loadReviews()
    })
</script>

<main>
    <h1>Revues</h1>
    <h2>Nouveau</h2>
    <div id="new-review-form">
        <label for="newReviewName">Nom</label>
        <input id="newReviewName" type="text" bind:value={newReview.name}>
        <label for="newReviewPassword">Mot de passe</label>
        <input type="text" bind:value={newReview.password}>
        <button on:click={handleNewReview}>Valider</button>
    </div>

    <div id="reviews">
        <div id="reviews-header">
            <div class="cell">Nom</div>
            <div class="cell">Mot de passe</div>
            <div class="cell">Photos</div>
            <div class="actions">Actions</div>
        </div>
        {#each reviews as review}
            <div class="review">
                <div class="cell">{review.name}</div>
                <div class="cell">{review.password}</div>
                <div class="cell">{review.pictures.length}</div>
                <div class="actions">
                    <button class="btn-error" on:click={() => handleView(review.name)}><img class="icon" src={View}
                                                                                            alt=""></button>
                    <button class="btn-error" on:click={() => handleRefresh(review.name)}><img class="icon" src={Renew}
                                                                                               alt=""></button>
                    <button class="btn-error" on:click={() => handleLoadEvents(review.name)}><img class="icon"
                                                                                                  src={List}
                                                                                                  alt=""></button>
                    <button class="btn-error" on:click={() => handleDelete(review.name)}><img class="icon"
                                                                                              src={Trashcan} alt="">
                    </button>
                </div>
            </div>
        {/each}
    </div>

    {#if showEvents}
        <div id="events">
            <h2>Evenements</h2>
            <div id="events-header" class="events-row">
                <div class="events-cell">Fichier</div>
                <div class="events-cell">Miniature</div>
                <div class="events-cell">BDD</div>
            </div>
            {#each Object.keys(events) as fileName}
                <div class="events-row">
                    <div class="events-cell">{fileName}</div>
                    <Status status={events[fileName].half}/>
                    <Status status={events[fileName].db}/>
                </div>
            {/each}

            <h2>Archives</h2>
            <Status status={archivesDone}/>
        </div>
    {/if}
</main>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  main {
    @include f-p-2;
    width: 80vw;
    position: absolute;
    top: 10vh;
    left: 10vw;

    h1, h2 {
      @include f-p-b;
    }

    #new-review-form {
      @include f-p-2;
      display: flex;
      flex-direction: column;

      input, label, button {
        width: 100%;
        margin-top: 10px;
      }

      label, input {
        font-size: 1.5rem;
      }

      button {
        @include f-p-b;
        font-size: 2em;
        outline: none !important;
        border: none !important;
        padding: 5px;
        box-shadow: 5px 5px 5px rgba(0, 0, 0, 0.2);
        border-radius: 5px;
        background-color: $subcolor;
        color: white;

        &:hover {
          cursor: pointer;
        }
      }
    }

    #reviews {
      width: 100%;

      #reviews-header {
        @include f-p-b;
        margin-top: 50px;
      }

      #reviews-header, .review {
        height: 30px;
        display: flex;

        .cell {
          width: 30%;
        }

        .actions {
          width: 10%;
        }
      }

      .review {
        .cell {
          width: 30%;
        }

        .actions {
          display: flex;
          width: 10%;
          align-items: center;
          justify-content: space-between;
        }


        button {
          border: none;
          outline: none;
          background: none;
          color: black;
          display: flex;
          justify-content: center;
          align-items: center;
          height: 20px;

          img {
            height: 20px;
            width: 20px;
          }

          &:hover {
            cursor: pointer;
          }

          &.btn-error {
            color: red;
          }
        }
      }
    }

    #events {
      width: 100%;

      #events-header {
        @include f-p-b;
        font-size: 1.3em;
        border-bottom: 1px solid $subcolor
      }

      .events-row {
        display: flex;
        justify-content: space-between;
        margin-top: 5px;

        .events-cell {
          width: 30%;
        }
      }
    }
  }
</style>