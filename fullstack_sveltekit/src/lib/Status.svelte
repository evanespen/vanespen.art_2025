<script>
    import Done from '$lib/svgs/checkmark--filled.svg?component';
    import Error from '$lib/svgs/error.svg?component';
    import Pending from '$lib/svgs/pending.svg?component';

    export let status;
    export let size = 20;
</script>

<main>
    {#if status === undefined}
        <Pending {size} fill="#0288D1"/>
    {:else if status === true}
        <Done {size} fill="#689F38"/>
    {:else if status === false}
        <Error {size} fill="#E64A19"/>
    {/if}
</main>

<style lang="scss">

</style>