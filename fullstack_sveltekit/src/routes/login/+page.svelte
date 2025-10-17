<script lang="ts">
    import {invalidateAll} from '$app/navigation';
    import {applyAction} from '$app/forms';

    /** @type {import('./$types').PageData} */
    export let data;

    /** @type {import('./$types').ActionData} */
    export let form;

    let invalid = false;

    async function handleSubmit() {
        const response = await fetch(this.action, {
            method: 'POST',
            body: new FormData(this)
        });

        /** @type {import('@sveltejs/kit').ActionResult} */
        const result = await response.json();

        if (result.type === 'success') {
            invalid = false;
            await invalidateAll();
        } else {
            invalid = true;
        }

        applyAction(result);
    }

</script>

<main id="login-page">
    <form method="POST" on:submit|preventDefault={handleSubmit}>
        {#if invalid}
            <div class="error">Erreur</div>
        {/if}
        <label for="username">Nom Utilisateur</label>
        <input id="username" name="username" type="text" value={form?.username ?? ''}>

        <label for="password">Mot de passe</label>
        <input id="password" name="password" type="password">
        <button>Connexion</button>
    </form>
</main>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  #login-page {
    position: absolute;
    top: 7vh;
    left: 10vw;
    width: 80vw;
    height: 80vh;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    form {
      display: flex;
      flex-direction: column;
      justify-content: center;
      width: 50%;

      .error {
        @include f-p-2;
        border: 1px solid red;
        border-radius: 5px;
        color: red;
        font-size: 2em;
        text-align: center;
        margin-bottom: 10px;
      }

      label {
        @include f-p-2;
        font-size: 1.5em;
      }

      input, button {
        margin-bottom: 20px;
        border: 1px solid $text;
        background-color: white !important;
        font-size: 1.5em;
        border-radius: 5px;
        box-shadow: 5px 5px 5px rgba(0, 0, 0, .2);
      }

      button {
        @include f-p-2;
        font-size: 1.5em;
      }
    }
  }
</style>
