<script>
    import SideTitles from "$lib/SideTitles.svelte";
    import ReviewPicture from "$lib/ReviewPicture.svelte";
    import {onMount} from "svelte";

    /** @type {import('./$types').PageData} */
    export let data;
    const review = data.review;

    let form;
    let authorized = data.authorized;

    function checkAuth(evt) {
        evt.preventDefault();
        form.incorrect = false;
        const formData = new FormData(form);

        fetch('?/checkAuth', {
            method: 'POST',
            body: formData
        }).then(res => {
            res.json().then(res => {
                form.incorrect = !res.data.auth;
                authorized = res.data.auth;
            })
        });
    }

    function downloadAll() {
        data.archives.forEach(archiveName => {
            const src = `/api/review-pictures/${review.name}/${archiveName}`;
            const a = document.createElement('a');
            a.href = src;
            a.download = src.split('/').pop();
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        })
    }

    onMount(() => {
    })
</script>

<main>

    <SideTitles/>

    <div id="review-header">
        <h2>Revue de la séance "<i>{data.review.name}</i>"</h2>
        <h5>({data.review.pictures.length} photos)</h5>
    </div>

    {#if authorized}
        <div id="download-button">
            <button on:click={downloadAll}>Tout télécharger</button>
        </div>

        <div id="review-pictures">
            {#each data.review.pictures as picture}
                <ReviewPicture {picture} {review}/>
            {/each}
        </div>
    {:else}
        <div id="review-login">
            <div class="review-login-error">
                <p>Vous n'êtes pas autorisé à télécharger les photos de la séance.</p>
                <p>Veuillez entrer le mot de passe pour accéder aux photos de la séance.</p>
            </div>

            {#if form?.incorrect}
                <div class="review-login-error">
                    <p class="error">Mot de passe incorrect !</p>
                </div>
            {/if}

            <form id="review-login-form" bind:this={form} on:submit={checkAuth}>
                <input name="password" type="text" placeholder="Mot de passe..." required>
                <button type="submit">Valider</button>
            </form>
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
    padding-bottom: 10vh !important;

    button {
      @include f-p-b;
      width: 100%;
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

    #review-header {
      width: 100%;
      border-bottom: 1px solid $text;
      height: 3em;
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;

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

    #review-login {
      .review-login-error {
        background-color: #E64A19;
        color: white;
        padding: 2px 10px;
        border-radius: 5px;
        margin-bottom: 10px;
      }

      #review-login-form {
        margin-top: 10px;
        width: 100%;
        display: flex;
        flex-direction: column;

        input {
          margin-bottom: 10px;
        }
      }
    }

    #download-button {
      margin-top: 2vh;
    }

    #review-pictures {
      margin-top: 2vh;
      width: 100%;
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      grid-row-gap: 5vh;
    }
  }
</style>