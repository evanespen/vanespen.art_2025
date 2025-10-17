<script>
    import {onMount} from 'svelte'

    import lax from 'lax.js'

    let mainCT, landingImg, landingTitle, landingSubtitle, landingCTA, homeContent, homeContentImages
    let oldPos = 0


    // function animate (e) {
    //     let pos = (window.scrollY / window.innerHeight) * 100
    //
    //     if (pos > oldPos) {
    //         if (pos <  100) window.scrollTo(0, window.innerHeight)
    //     } else {
    //         if (pos > 80 && pos < 100) window.scrollTo(0, 0)
    //     }
    //
    //     if (pos > 20) landingCTA.style.visibility = 'hidden'
    //     if (pos < 20) landingCTA.style.visibility = 'visible'
    //
    //     landingImg.style.height = `${90 - pos}vh`
    //     landingImg.style.width = `${95 - pos}vw`
    //
    //     if (pos > 20) {
    //         landingTitle.style.transform = `rotate(-90deg)`
    //         landingTitle.style.marginLeft = '-90vw'
    //         landingTitle.style.color = '#211a17'
    //
    //         landingSubtitle.style.transform = `rotate(90deg)`
    //         landingSubtitle.style.marginRight = '-90vw'
    //         landingSubtitle.style.color = '#211a17'
    //
    //     } else {
    //         landingTitle.style.transform = `rotate(0deg)`
    //         landingTitle.style.marginLeft = '0'
    //         landingTitle.style.color = '#fff8ee'
    //
    //         landingSubtitle.style.transform = `rotate(0deg)`
    //         landingSubtitle.style.marginRight = '0'
    //         landingSubtitle.style.color = '#fff8ee'
    //     }
    //
    //     if (pos > 90) {
    //         landingImg.style.height = `0vh`
    //         landingImg.style.width = `0vw`
    //     }
    //
    //     oldPos = pos
    // }

    onMount(() => {
        // window.addEventListener('scroll', animate)

        lax.init()

        lax.addDriver('scrollY', function () {
            return window.scrollY
        })

        lax.addDriver('posY', function () {
            const posY = (window.scrollY / window.innerHeight) * 100

            if (landingCTA) {
                if (posY > 20) landingCTA.style.display = 'none'
                if (posY < 20) landingCTA.style.display = 'block'
            }

            if (landingTitle) {
                if (posY > 90) {
                    // landingTitle.style.color = '#211a17'
                    // landingSubtitle.style.color = '#211a17'
                    landingTitle.classList.add('scrolled')
                    landingSubtitle.classList.add('scrolled')
                } else {
                    // landingTitle.style.color = '#fff8ee'
                    // landingSubtitle.style.color = '#fff8ee'
                    landingTitle.classList.remove('scrolled')
                    landingSubtitle.classList.remove('scrolled')
                }
            }

            return posY > 100 ? 100 : posY
        })

        lax.addElements('#landing-title', {
            posY: {
                rotate: [
                    [0, 50, 90],
                    [0, -45, -90],
                ],
                translateX: [
                    [95, 100],
                    [0, '-screenWidth/2.2']
                ]
            }
        })

        lax.addElements('#landing-subtitle', {
            posY: {
                rotate: [
                    [0, 50, 90],
                    [0, 45, 90]
                ],
                translateX: [
                    [95, 100],
                    [0, 'screenWidth/2.2']
                ]
            }
        })

        lax.addElements('#landing-image', {
            posY: {
                scale: [
                    [0, 100],
                    [1, 0]
                ]
            }
        })

        // lax.addElements('.home-content-row-images', {
        //   scrollY: {
        // 	translateX: [
        // 	  ['elInY', 'elCenterY', 'elOutY'],
        // 	  [0, 500, 1000],
        // 	],
        //   }
        // })

        // homeContentImages.querySelectorAll('img').forEach(el => {
        //   console.log(el.src)
        //   el.style.marginLeft = `${el.dataset.scrollSpeed * 10}px`

        //   lax.addElements(`img[src=${el.src}]`, {
        // 	scrollY: {
        //     translateY: [
        //       ["elInY", "elOutY"],
        //       [0, `elCenterY/${el.dataset.scrollSpeed}`],
        //     ]
        // 	}
        //   })
        // })

        document.querySelectorAll('.home-content-row-images').forEach(hcri => {
            hcri.querySelectorAll('img').forEach(i => {
                const rect = i.getBoundingClientRect()
                const shadowRect = document.querySelector('#shadow-image-position').getBoundingClientRect()

                if (rect.left > shadowRect.right) {
                    i.style.transformOrigin = 'right'
                }

                i.addEventListener('mouseover', (e) => {
                    if (i.dataset.animated == true) return
                    i.dataset.animated = true


                    let translateYBy = -(rect.top - shadowRect.top)


                    let translateXBy
                    if (rect.left < shadowRect.right) translateXBy = -(rect.left - shadowRect.left) + 200
                    else translateXBy = (window.innerWidth / 2 - rect.left) + 200


                    i.style.transform = `scale(2)` //`translateY(${translateYBy}px) translateX(${translateXBy}px) scale(2)`
                    return
                })
                i.addEventListener('mouseleave', (e) => {
                    i.style.transform = 'none'
                    i.dataset.animated = false
                })
            })
        })
    })
</script>

<svelte:head>
    <title>Evrard Van Espen, photographie</title>
    <meta name="description" content="Evrard Van Espen, photographie"/>
</svelte:head>

<main id="home" bind:this={mainCT}>

    <div id="shadow-image-position"></div>

    <div id="landing">
        <!--        <img bind:this={landingImg}-->
        <!--             id="landing-image"-->
        <!--             src='https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/195a93cd-5675-43dc-b4c5-d80c14448f00/public'/>-->
        <div bind:this={landingImg}
             id="landing-image"></div>
        <div bind:this={landingTitle}
             id="landing-title">EVRARD VAN ESPEN
        </div>
        <div bind:this={landingSubtitle}
             id="landing-subtitle">Photographie
        </div>
    </div>

    <div id="landing-cta-container">
        <img bind:this={landingCTA}
             on:click={() => {scrollTo(0, window.innerHeight * 2)}}
             id="landing-cta"
             src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/a52bb9b9-0f2a-4c94-dd10-fc6d37987400/public"/>
    </div>

    <div id="home-content" bind:this={homeContent}>

        <div class="home-content-row">
            <div class="home-content-row-title"><p>Portraits</p></div>
            <div class="home-content-row-images">
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/e6c0385e-1504-4435-4f82-3037103ea000/public"/>
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/a7fda2b1-49df-4e17-d7c8-645285d27d00/public"/>
            </div>
        </div>

        <div class="home-content-row right">
            <div class="home-content-row-title"><p>Sports Automobiles</p></div>
            <div class="home-content-row-images">
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/f02099cc-4ffd-49b5-7e31-8e84b3de2a00/public"/>
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/fe2c9dd9-980d-4e92-3ba7-af9cbad64400/public"/>
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/a61712f4-a52d-4f42-fed7-aa1d07503000/public"/>
            </div>
        </div>

        <div class="home-content-row">
            <div class="home-content-row-title"><p>Photo Animalière</p></div>
            <div class="home-content-row-images">
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/9cf7e2b4-7b27-4de1-250a-f94bd092c000/public"/>
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/e6485aab-1416-4358-0f34-0806a8c97600/public"/>
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/8d4b6151-ce3a-4c2f-f2f8-d9c19ac98100/public"/>
                <img src="https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/a5142786-781d-44a1-3629-e27fce861f00/public"/>
            </div>
        </div>


    </div>

</main>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  #home {
    #shadow-image-position {
      /* border: 5px solid red; */
      position: fixed;
      top: 5vh;
      height: 90vh;
      width: 50vw;
    }

    #landing {
      pointer-events: none;
      position: fixed;
      top: 5vh;
      left: 0;
      height: 95vh;
      width: 100vw;
      z-index: 1000;

      display: flex;
      justify-content: center;
      align-items: center;
      flex-direction: column;
      color: $background;

      #landing-image {
        background-image: url(https://imagedelivery.net/SBTKyitl9TX77hzo56nCqA/195a93cd-5675-43dc-b4c5-d80c14448f00/public);
        background-size: cover;
        background-position: center;
        position: absolute;
        height: 90vh;
        width: 95vw;
        transition: .5s;
      }

      #landing-title {
        @include f-h;
        width: 100vh;
        display: flex;
        justify-content: center;

        font-size: 6em;
        z-index: 100;

        transition: .5s;
        color: #fff8ee;

        @media (max-width: 800px) {
          font-size: 11vw;
        }
      }

      #landing-subtitle {
        @include f-h;

        width: 100vh;
        display: flex;
        justify-content: center;

        font-size: 4em;
        font-variant: all-small-caps;
        margin-top: -20px;
        z-index: 100;

        transition: .5s;
        color: #fff8ee;

        @media (max-width: 800px) {
          font-size: 8vw;
        }
      }
    }

    #landing-cta-container {
      width: 100vw;
      position: fixed;
      top: 80vh;
      z-index: 1000000;
      display: flex;
      justify-content: center;
      align-items: center;

      #landing-cta {
        height: 10vh;
        width: 10vh;
        filter: invert(1);
        animation: pulse 2s infinite;
        transition: .5s;
        display: flex;

        &:hover {
          cursor: pointer;
        }
      }
    }

    #home-content {
      position: absolute;
      top: 200vh;
      left: 10%;
      height: 100vh;
      width: 80%;
      z-index: 1000;

      /* #home-content-images { */
      /* 	display: flex; */
      /* 	flex-wrap: wrap; */
      /* 	justify-content: space-between; */

      /* 	img { */
      /* 	  max-width: 40%; */
      /* 	  margin: 20px; */
      /* 	  margin-bottom: 100px; */
      /* 	} */
      /* } */

      .home-content-row {
        margin-bottom: 100px;
        display: flex;
        flex-direction: column;


        &.right {
          .home-content-row-title {
            display: flex;
            flex-direction: column;
            align-items: flex-end;
          }

          .home-content-row-images {
            justify-content: flex-end;

            img {
              margin-right: 0;
              margin-left: 10px;
            }
          }
        }

        .home-content-row-title {
          @include f-h;
          font-size: 2em;
          display: flex;
          align-items: center;
          border-bottom: 1px solid $text;
          margin-bottom: 20px;

          p {
            padding: 0;
            margin: 0;
            margin-bottom: -10px;
          }
        }

        .home-content-row-images {
          display: flex;

          @media (max-width: 800px) {
            flex-direction: row;
            justify-content: space-between;
            flex-wrap: wrap;
            width: 100%;
          }

          img {
            max-height: 30vh;
            margin-right: 10px;
            transition: .5s;
            transform-origin: left top;

            @media (max-width: 800px) {
              margin: 0 0 10px;
              max-width: 100%;
            }

            /* &:hover { */
            /*   max-height: 80vh; */
            /*   transform: scale(2); */
            /* } */
          }
        }
      }
    }
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.2);
    }
  }

  :global(.scrolled) {
    color: #211a17 !important;

    @media (max-width: 800px) {
      opacity: 0 !important;
    }
  }
</style>
