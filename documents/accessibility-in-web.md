---
title: "ITP: Accessibility in Web"
tags:
  - Assignment
  - ITP
categories:
  - Essay
date: 2018-10-09 13:21:10
---

### Intro

Internet is a global web that connects people. It provides services or information for everyone in the world. So a web should be well prepared for anyone to get access to it. A standard for a *accessible* page is also standardized by MDN group (Mozilla) in its [Accessibility Page](https://developer.mozilla.org/en-US/docs/Web/Accessibility). By using some web design or html attributions, we can make the webpage clear and easy for most of people.

In the **attempt for accessibility**, I mainly focus on two projects - [Nature of Code Book Page](https://natureofcode.com/book/) and *a School PWA Project* (since it is not released yet, so I just hide the URL to that webapp). During the research, I find it interesting that accessibility is **hard to test** if we see it as our origin aspect. That means, to test the accessibility for a web, we need to *pretend* to stand on the aspect of people who actually need this function. And I well talk about my research about *nature of code book* in the following paragraph:

### Test with Color Contrast

Firstly, by using [Color Contrast Analyser](https://www.paciellogroup.com/resources/contrastanalyser/) (since we cannot even tell the difference of color contrast with our eyes usually), I test the primary(#D82F5C), secondary(#61D2EA) and highlight(#0019B2) colors of *nature of code book* page.

Following image is the result of testing of primary color(#D82F5C):

![Accessilibity 1](https://yuuno.cc/images/accessibility-1.png)

We can see from the result, the color is perfect besides non-title regular test. But it is fine since the page only use this color with **bold text** or behind **title** division.

And the two other colors: secondary color(#61D2EA) and highlight color(#0019B2):

![Accessilibity 2](https://yuuno.cc/images/accessibility-2.png)

![Accessilibity 3](https://yuuno.cc/images/accessibility-3.png)

I find it very interesting that the button color(#61D2EA) does not pass even one test in color contract, but if user interacts with the web by using **tab** or **VoiceOver** control, the secondary color transits to the highlight color(#0019B2) which totally passes all the tests.

So in this case, the button is friendly for *Accessibility Tools* user to get full support on the page.

### Test with Tab

Another test is using *Tab* key to test the web.

In this test, I dropped mouse and only use *tab* and *space/enter* to interact with the web. And find an issue that [some hidden buttons in the page might confusing users](https://github.com/shiffman/natureofcode.com/pull/30).

So I made a pull request to correct this issue. To show directly the button if user uses *tab* or other *accessibility tool* to surf the page.
