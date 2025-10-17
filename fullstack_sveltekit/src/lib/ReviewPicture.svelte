<script lang="ts">
    import ThumbsUp from '$lib/svgs/thumbs-up.svg?component';
    import ThumbsUpFilled from '$lib/svgs/thumbs-up--filled.svg?component';
    import ThumbsDown from '$lib/svgs/thumbs-down.svg?component';
    import ThumbsDownFilled from '$lib/svgs/thumbs-down--filled.svg?component'
    import Send from '$lib/svgs/send.svg?component';
    import Download from '$lib/svgs/download.svg?component';

    export let picture;
    export let review;

    const pictureFname = picture.path.split('/').at(-1);
    const src = `/api/review-pictures/${review.name}/${pictureFname}`;

    let comment = picture.comment;

    async function setStatus(status) {
        picture.status = status;
        console.log(picture.status);
        console.log(src)

        const response = await fetch(src, {
            method: 'PUT',
            body: JSON.stringify({
                action: 'setStatus',
                value: picture.status,
            })
        });
    }

    async function setComment() {
        picture.comment = comment;

        const response = await fetch(src, {
            method: 'PUT',
            body: JSON.stringify({
                action: 'setComment',
                value: picture.comment,
            })
        });
    }

    function download() {
        const a = document.createElement('a');
        a.href = src;
        a.download = src.split('/').pop();
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
    }

</script>

<main>
    <img src={`${src}?type=half`}>
    <div id="controls">
        <div id="buttons">

            {#if picture.status === -1}
                <button on:click={() => setStatus(1)}>
                    <ThumbsUp/>
                </button>
                <button on:click={() => setStatus(0)}>
                    <ThumbsDownFilled fill="#E64A19"/>
                </button>
            {/if}
            {#if picture.status === 0}
                <button on:click={() => setStatus(1)}>
                    <ThumbsUp/>
                </button>
                <button on:click={() => setStatus(-1)}>
                    <ThumbsDown/>
                </button>
            {/if}
            {#if picture.status === 1}
                <button on:click={() => setStatus(0)}>
                    <ThumbsUpFilled fill="#689F38"/>
                </button>
                <button on:click={() => setStatus(-1)}>
                    <ThumbsDown/>
                </button>
            {/if}

            <button on:click={download}>
                <Download/>
            </button>

        </div>
        <div id="comment">
            <textarea bind:value={comment} placeholder="Commentaire..."></textarea>
            <button on:click={() => setComment()}>
                <Send/>
            </button>
        </div>
    </div>
</main>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  main {
    border-radius: 0 0 10px 10px;
    background-color: white;
    box-shadow: 5px 5px 20px rgba(0, 0, 0, 0.2);
    width: 25vw;
    align-self: start;

    @media (max-width: 800px) {
      width: 100%;
    }


    img {
      width: 100%;
    }

    #controls {
      height: 100px;
      display: flex;
      justify-content: space-between;
      padding: 0px 10px;

      #buttons {
        display: flex;
        align-items: center;
        justify-content: space-around;
        width: 25%;

        button {
          height: 32px;
          border: none;
          outline: none;
          background: none;

          &:hover {
            cursor: pointer;
          }

          img {
            width: 100%;
            height: 100%;
          }
        }
      }

      #comment {
        width: 70%;
        display: flex;
        justify-content: center;
        align-items: center;

        button {
          width: 32px;
          border: none;
          outline: none;
          background: none;

          &:hover {
            cursor: pointer;
          }

          img {
            width: 100%;
            height: 100%;
          }
        }

        textarea {
          width: 100%;
          height: 80%;
          resize: none;
          font-family: sans-serif;
        }
      }
    }
  }
</style>