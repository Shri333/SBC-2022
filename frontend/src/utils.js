function toInt(value) {
  const res = Number(value);
  if (isNaN(res))
    return "NaN";
  return res;
}

export { toInt };
