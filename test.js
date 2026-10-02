/**
 * @param {number} x
 * @return {boolean}
 */
var isPalindrome = function (x) {
  console.log(x);
  const xString = x.toString();
  const isEven = xString.length % 2 === 0;
  let isPalindromeTrue = false;
  if (x < 0) {
    return false;
  }
  if (xString.length === 1) {
    return true;
  }

  console.log("isEven", isEven);
  if (isEven) {
    const firstHalf = xString.slice(0, xString.length / 2);
    const secondHalf = xString
      .slice(xString.length / 2, xString.length)
      .split("")
      .reverse()
      .join("");
    console.log("firstHalf", firstHalf);
    console.log("secondHalf", secondHalf);
    isPalindromeTrue = firstHalf === secondHalf;
  } else {
    const firstHalf = xString.slice(0, Math.floor(xString.length / 2));
    const secondHalf =
      xString.slice(Math.floor(xString.length / 2) + 1, xString.length).at(-1) +
      xString
        .slice(Math.floor(xString.length / 2) + 1, xString.length - 1)
        .split("")
        .reverse()
        .join("");
    console.log("firstHalf", firstHalf);
    console.log("secondHalf", secondHalf);
    isPalindromeTrue = firstHalf === secondHalf;
  }
  console.log("isPalindromeTrue", isPalindromeTrue);

  return isPalindromeTrue;
};

isPalindrome(1001);
